package k8s

import (
	"context"
	"errors"
	"fmt"
	"gim/pkg/logger"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sort"
	"strings"

	"google.golang.org/grpc/resolver"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// 实现k8s地址解析，根据k8s的service的endpoints解析 比如，k8s:///namespace.server:port
func init() {
	resolver.Register(NewK8sBuilder())
}

func GetK8STarget(namespace, server, port string) string {
	return fmt.Sprintf("k8s:///%s.%s:%s", namespace, server, port)
}

type k8sBuilder struct{}

func NewK8sBuilder() resolver.Builder {
	return &k8sBuilder{}
}

func (b *k8sBuilder) Build(target resolver.Target, clientConn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	return newK8sResolver(target, clientConn)
}

func (b *k8sBuilder) Scheme() string {
	return "k8s"
}

// k8sResolver k8s地址解析器
type k8sResolver struct {
	clientConn resolver.ClientConn
	endpoint   string
	ips        []string
	port       string
	cancel     context.CancelFunc
}

func getK8sClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return k8sClient, nil
}

func newK8sResolver(target resolver.Target, clientConn resolver.ClientConn) (*k8sResolver, error) {
	namespace, service, port, err := resolveEndpoint(target.Endpoint)
	if err != nil {
		return nil, err
	}

	k8sClient, err := getK8sClient()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	k8sResolver := &k8sResolver{
		clientConn: clientConn,
		endpoint:   target.Endpoint,
		port:       port,
		cancel:     cancel,
	}

	// 监听变化
	go func() {
		events, err := k8sClient.CoreV1().Endpoints(namespace).Watch(context.TODO(), metav1.ListOptions{
			LabelSelector: fields.OneTermEqualSelector("app", service).String(),
		})
		if err != nil {
			logger.Logger.Error("k8s resolver error", zap.String("endpoint", target.Endpoint), zap.Error(err))
			panic(err)
		}

		for {
			select {
			case <-ctx.Done():
				return
			case event := <-events.ResultChan():
				endpoints, ok := event.Object.(*v1.Endpoints)
				if !ok {
					logger.Logger.Warn("k8s resolver error", zap.String("endpoint", target.Endpoint), zap.Any("event", event))
					continue
				}
				ips, err := getIpsFromEndpoint(endpoints)
				if err != nil {
					logger.Logger.Warn("k8s resolver error", zap.String("endpoint", target.Endpoint), zap.Error(err))
					continue
				}
				k8sResolver.updateState(ips)
			}
		}
	}()
	return k8sResolver, nil
}

func (r *k8sResolver) ResolveNow(opt resolver.ResolveNowOptions) {
	logger.Logger.Info("k8s resolver resolveNow", zap.String("endpoint", r.endpoint))
}

func (r *k8sResolver) Close() {
	r.cancel()
}

// updateState 更新地址列表
func (r *k8sResolver) updateState(newIPs []string) {
	if isEqualIps(r.ips, newIPs) {
		return
	}
	r.ips = newIPs

	addresses := make([]resolver.Address, 0, len(r.ips))
	for _, v := range r.ips {
		addresses = append(addresses, resolver.Address{
			Addr: v + ":" + r.port,
		})
	}
	state := resolver.State{
		Addresses: addresses,
	}
	logger.Logger.Info("k8s resolver updateState", zap.String("endpoint", r.endpoint), zap.Any("addresses", addresses))
	r.clientConn.UpdateState(state)
	return
}

// resolveEndpoint 对grpc的Endpoint进行解析，格式必须是：k8s:///namespace.server:port
func resolveEndpoint(endpoint string) (namespace string, service string, port string, err error) {
	namespaceAndServerPort := strings.Split(endpoint, ".")
	if len(namespaceAndServerPort) != 2 {
		err = errors.New("endpoint must is namespace.server:port")
		return
	}
	namespace = namespaceAndServerPort[0]
	serverAndPort := strings.Split(namespaceAndServerPort[1], ":")
	if len(serverAndPort) != 2 {
		err = errors.New("endpoint must is namespace.server:port")
		return
	}
	service = serverAndPort[0]
	port = serverAndPort[1]
	return
}

// isEqualIps 判断两个地址列表是否相等
func isEqualIps(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	sort.Strings(s1)
	sort.Strings(s2)
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// getIpsFromEndpoint 获取endpoints里面的IP列表
func getIpsFromEndpoint(endpoints *v1.Endpoints) ([]string, error) {
	if len(endpoints.Subsets) < 1 {
		return nil, errors.New("subsets length less than 1")
	}
	endpointAddresses := endpoints.Subsets[0].Addresses

	ips := make([]string, 0, len(endpoints.Subsets[0].Addresses))
	for _, v := range endpointAddresses {
		ips = append(ips, v.IP)
	}
	return ips, nil
}

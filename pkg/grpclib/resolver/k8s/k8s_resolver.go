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
	"time"

	"google.golang.org/grpc/resolver"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var k8sClientSet *kubernetes.Clientset

func GetK8sClient() (*kubernetes.Clientset, error) {
	if k8sClientSet == nil {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		k8sClientSet, err = kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}
	}
	return k8sClientSet, nil
}

// 实现k8s地址解析，根据k8s的service的endpoints解析 比如，k8s:///namespace.server:port
func init() {
	resolver.Register(&k8sBuilder{})
}

func GetK8STarget(namespace, server, port string) string {
	return fmt.Sprintf("k8s:///%s.%s:%s", namespace, server, port)
}

type k8sBuilder struct{}

func (b *k8sBuilder) Build(target resolver.Target, clientConn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	return newK8sResolver(target, clientConn)
}

func (b *k8sBuilder) Scheme() string {
	return "k8s"
}

// k8sResolver k8s地址解析器
type k8sResolver struct {
	log             *zap.Logger
	clientConn      resolver.ClientConn
	endpointsClient corev1.EndpointsInterface
	service         string

	cancel context.CancelFunc

	ips  []string
	port string
}

func newK8sResolver(target resolver.Target, clientConn resolver.ClientConn) (*k8sResolver, error) {
	log := logger.Logger.With(zap.String("target", target.Endpoint))
	log.Info("k8s resolver build")
	namespace, service, port, err := parseTarget(target)
	if err != nil {
		log.Error("k8s resolver error", zap.Error(err))
		return nil, err
	}

	k8sClient, err := GetK8sClient()
	if err != nil {
		log.Error("k8s resolver error", zap.Error(err))
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	client := k8sClient.CoreV1().Endpoints(namespace)
	k8sResolver := &k8sResolver{
		log:             log,
		clientConn:      clientConn,
		endpointsClient: client,
		service:         service,
		cancel:          cancel,
		port:            port,
	}
	err = k8sResolver.updateState(true)
	if err != nil {
		log.Error("k8s resolver error", zap.Error(err))
		return nil, err
	}

	ticker := time.NewTicker(time.Second)
	// 监听变化
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = k8sResolver.updateState(false)
			}
		}
	}()
	return k8sResolver, nil
}

// ResolveNow grpc感知到连接异常，会做通知，观察日志得知
func (r *k8sResolver) ResolveNow(opt resolver.ResolveNowOptions) {
	r.log.Info("k8s resolver resolveNow")
}

func (r *k8sResolver) Close() {
	r.log.Info("k8s resolver close")
	r.cancel()
}

// updateState 更新地址列表
func (r *k8sResolver) updateState(isFromNew bool) error {
	endpoints, err := r.endpointsClient.Get(context.TODO(), r.service, metav1.GetOptions{})
	if err != nil {
		r.log.Error("k8s resolver error", zap.Error(err))
		return err
	}
	newIPs := getIPs(endpoints)
	if len(newIPs) == 0 {
		return nil
	}
	if !isFromNew && isEqualIPs(r.ips, newIPs) {
		return nil
	}
	r.ips = newIPs

	addresses := make([]resolver.Address, 0, len(r.ips))
	for _, ip := range r.ips {
		addresses = append(addresses, resolver.Address{
			Addr: ip + ":" + r.port,
		})
	}
	state := resolver.State{
		Addresses: addresses,
	}
	r.log.Info("k8s resolver updateState", zap.Bool("is_from_new", isFromNew), zap.Any("service", r.service), zap.Any("addresses", addresses))
	// 这里地址数量不能为0，为0会返回错误
	err = r.clientConn.UpdateState(state)
	if err != nil {
		r.log.Error("k8s resolver error", zap.Error(err))
		return err
	}
	return nil
}

// parseTarget 对grpc的Endpoint进行解析，格式必须是：k8s:///namespace.server:port
func parseTarget(target resolver.Target) (namespace string, service string, port string, err error) {
	namespaceAndServerPort := strings.Split(target.Endpoint, ".")
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

// isEqualIPs 判断两个地址列表是否相等
func isEqualIPs(s1, s2 []string) bool {
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

// getIPs 获取EndpointSlice里面的IP列表
func getIPs(endpoints *v1.Endpoints) []string {
	ips := make([]string, 0, 10)
	if len(endpoints.Subsets) <= 0 {
		return ips
	}

	for _, address := range endpoints.Subsets[0].Addresses {
		ips = append(ips, address.IP)
	}
	return ips
}

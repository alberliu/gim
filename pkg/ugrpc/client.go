package ugrpc

import (
	"fmt"

	"github.com/sercand/kuberesolver/v6"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gim/config"
)

func init() {
	kuberesolver.RegisterInClusterWithSchema("k8s")
}

func NewClient(target string) *grpc.ClientConn {
	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(clientInterceptor),
		WithDefaultServiceConfig(),
	)
	if err != nil {
		panic(err)
	}
	return conn
}

func WithDefaultServiceConfig() grpc.DialOption {
	return grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`)
}

func GetTarget(server string) string {
	switch config.ENV {
	case config.EnvCompose:
		return fmt.Sprintf("dns:///%s:8000", server)
	case config.EnvK8s:
		return fmt.Sprintf("k8s:///%s:8000", server)
	default:
		return ""
	}
}

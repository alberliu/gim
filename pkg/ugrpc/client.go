package ugrpc

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func interceptor(ctx context.Context, method string, request, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, request, reply, cc, opts...)

	md, _ := metadata.FromOutgoingContext(ctx)
	slog.Debug("client interceptor", "method", method, "metadata", md, "request", request, "reply", reply, "error", err)
	return err
}

func NewClient(target string) *grpc.ClientConn {
	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(interceptor),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
	if err != nil {
		panic(err)
	}
	return conn
}

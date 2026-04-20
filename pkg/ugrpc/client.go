package ugrpc

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func clientInterceptor(ctx context.Context, method string, request, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	err := invoker(ctx, method, request, reply, cc, opts...)

	md, _ := metadata.FromOutgoingContext(ctx)
	slog.Debug("client interceptor", "method", method, "metadata", md, "request", request, "reply", reply, "error", err)
	return err
}

func NewClient(target string) *grpc.ClientConn {
	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(clientInterceptor),
		grpc.WithDefaultServiceConfig(`{
            "loadBalancingConfig": [{"round_robin":{}}]
        }`),
	)
	if err != nil {
		panic(err)
	}
	return conn
}

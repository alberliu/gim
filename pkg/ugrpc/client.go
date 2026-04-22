package ugrpc

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"gim/pkg/logger"
	"gim/pkg/md"
)

func clientInterceptor(ctx context.Context, method string, request, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	ctx = md.WithRequestID(ctx, md.GetRequestID(ctx))
	err := invoker(ctx, method, request, reply, cc, opts...)

	md, _ := metadata.FromOutgoingContext(ctx)
	slog.DebugContext(ctx, "client interceptor", "method", method, "metadata", md, "request", request, "reply", reply, logger.Error(err))
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

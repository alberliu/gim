package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"gim/pkg/protocol/pb/connectpb"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/protocol/pb/userpb"
)

const EnvLocal = "local"

var ENV = os.Getenv("ENV")

var builders = map[string]Builder{
	"local":   &localBuilder{},
	"compose": &composeBuilder{},
	"k8s":     &k8sBuilder{},
}

var Config Configuration

type Builder interface {
	Build() Configuration
}

type Configuration struct {
	LogLevel slog.Level
	LogFile  func(server string) string

	MySQL                string
	RedisHost            string
	RedisPassword        string
	PushRoomSubscribeNum int
	PushAllSubscribeNum  int

	ConnectLocalAddr     string
	ConnectRPCListenAddr string
	ConnectTCPListenAddr string
	ConnectWSListenAddr  string

	LogicRPCListenAddr string
	UserRPCListenAddr  string
	FileHTTPListenAddr string

	ConnectIntClientBuilder func() connectpb.ConnectIntServiceClient
	DeviceIntClientBuilder  func() logicpb.DeviceIntServiceClient
	MessageIntClientBuilder func() logicpb.MessageIntServiceClient
	RoomIntClientBuilder    func() logicpb.RoomIntServiceClient
	UserIntClientBuilder    func() userpb.UserIntServiceClient
}

func init() {
	builder, ok := builders[ENV]
	if !ok {
		builder = new(localBuilder)
	}
	Config = builder.Build()
}

func interceptor(ctx context.Context, method string, request, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, request, reply, cc, opts...)

	md, _ := metadata.FromOutgoingContext(ctx)
	slog.Debug("client interceptor", "method", method, "metadata", md, "request", request, "reply", reply, "error", err)
	return err
}

func newGrpcClient(target, loadBalance string) *grpc.ClientConn {
	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(interceptor),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, loadBalance)))
	if err != nil {
		panic(err)
	}
	return conn
}

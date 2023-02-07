package config

import (
	"context"
	"fmt"
	"gim/pkg/grpclib/picker"
	"gim/pkg/logger"
	"gim/pkg/protocol/pb"

	"go.uber.org/zap"

	_ "gim/pkg/grpclib/resolver/addrs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

type defaultBuilder struct{}

func (*defaultBuilder) Build() Configuration {
	logger.Level = zap.DebugLevel
	logger.Target = logger.Console

	return Configuration{
		MySQL:                "root:gim123456@tcp(111.229.238.28:3306)/gim?charset=utf8&parseTime=true",
		RedisHost:            "111.229.238.28:6379",
		RedisPassword:        "alber123456",
		PushRoomSubscribeNum: 100,
		PushAllSubscribeNum:  100,

		ConnectLocalAddr:     "127.0.0.1:8000",
		ConnectRPCListenAddr: ":8000",
		ConnectWSListenAddr:  ":8001",
		ConnectTCPListenAddr: ":8002",

		LogicRPCListenAddr:    ":8010",
		BusinessRPCListenAddr: ":8020",
		FileHTTPListenAddr:    "8030",

		ConnectIntClientBuilder: func() pb.ConnectIntClient {
			conn, err := grpc.DialContext(context.TODO(), "addrs:///127.0.0.1:8000", grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
				grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, picker.AddrPickerName)))
			if err != nil {
				panic(err)
			}
			return pb.NewConnectIntClient(conn)
		},
		LogicIntClientBuilder: func() pb.LogicIntClient {
			conn, err := grpc.DialContext(context.TODO(), "addrs:///docker.for.mac.host.internal:8010", grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
				grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
			if err != nil {
				panic(err)
			}
			return pb.NewLogicIntClient(conn)
		},
		BusinessIntClientBuilder: func() pb.BusinessIntClient {
			conn, err := grpc.DialContext(context.TODO(), "addrs:///127.0.0.1:8020", grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
				grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
			if err != nil {
				panic(err)
			}
			return pb.NewBusinessIntClient(conn)
		},
	}
}

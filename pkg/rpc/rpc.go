package rpc

import (
	"context"
	"fmt"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/pb"

	"google.golang.org/grpc/balancer/roundrobin"

	"google.golang.org/grpc"
)

var (
	LogicIntClient    pb.LogicIntClient
	ConnectIntClient  pb.ConnectIntClient
	BusinessIntClient pb.BusinessIntClient
)

func InitLogicIntClient(addr string) {
	conn, err := grpc.DialContext(context.TODO(), addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	LogicIntClient = pb.NewLogicIntClient(conn)
}

func InitConnectIntClient(addr string) {
	conn, err := grpc.DialContext(context.TODO(), addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, grpclib.AddrPickerName)))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	ConnectIntClient = pb.NewConnectIntClient(conn)
}

func InitBusinessIntClient(addr string) {
	conn, err := grpc.DialContext(context.TODO(), addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	BusinessIntClient = pb.NewBusinessIntClient(conn)
}

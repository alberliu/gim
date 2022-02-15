package rpc

import (
	"context"
	"fmt"
	"gim/config"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/pb"

	"google.golang.org/grpc/balancer/roundrobin"

	"google.golang.org/grpc"
)

var (
	logicIntClient    pb.LogicIntClient
	connectIntClient  pb.ConnectIntClient
	businessIntClient pb.BusinessIntClient
)

func GetLogicIntClient() pb.LogicIntClient {
	if logicIntClient == nil {
		initLogicIntClient()
	}
	return logicIntClient
}

func GetConnectIntClient() pb.ConnectIntClient {
	if connectIntClient == nil {
		initConnectIntClient()
	}
	return connectIntClient
}

func GetBusinessIntClient() pb.BusinessIntClient {
	if businessIntClient == nil {
		initBusinessIntClient()
	}
	return businessIntClient
}

func initLogicIntClient() {
	conn, err := grpc.DialContext(context.TODO(), config.RPCAddr.LogicRPCAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	logicIntClient = pb.NewLogicIntClient(conn)
}

func initConnectIntClient() {
	conn, err := grpc.DialContext(context.TODO(), config.RPCAddr.ConnectRPCAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, grpclib.AddrPickerName)))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	connectIntClient = pb.NewConnectIntClient(conn)
}

func initBusinessIntClient() {
	conn, err := grpc.DialContext(context.TODO(), config.RPCAddr.BusinessRPCAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	businessIntClient = pb.NewBusinessIntClient(conn)
}

package rpc

import (
	"context"
	"fmt"
	"gim/pkg/grpclib/picker"
	"gim/pkg/grpclib/resolver/k8s"
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
	conn, err := grpc.DialContext(context.TODO(), k8s.GetK8STarget("default", "logic", "8000"), grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
	if err != nil {
		panic(err)
	}
	logicIntClient = pb.NewLogicIntClient(conn)
}

func initConnectIntClient() {
	conn, err := grpc.DialContext(context.TODO(), k8s.GetK8STarget("default", "connect", "8000"), grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, picker.AddrPickerName)))
	if err != nil {
		panic(err)
	}
	connectIntClient = pb.NewConnectIntClient(conn)
}

func initBusinessIntClient() {
	conn, err := grpc.DialContext(context.TODO(), k8s.GetK8STarget("default", "business", "8000"), grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
	if err != nil {
		panic(err)
	}

	businessIntClient = pb.NewBusinessIntClient(conn)
}

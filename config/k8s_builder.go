package config

import (
	"context"
	"fmt"
	"gim/pkg/grpclib/picker"
	"gim/pkg/grpclib/resolver/k8s"
	"gim/pkg/k8sutil"
	"gim/pkg/logger"
	"gim/pkg/protocol/pb"
	"os"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type k8sBuilder struct{}

func (*k8sBuilder) Build() Configuration {
	const (
		RPCListenAddr = ":8000"
		RPCDialAddr   = "8000"
	)
	const namespace = "gim"

	k8sClient, err := k8sutil.GetK8sClient()
	if err != nil {
		panic(err)
	}
	configmap, err := k8sClient.CoreV1().ConfigMaps(namespace).Get(context.TODO(), "config", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	logger.Level = zap.DebugLevel
	logger.Target = logger.Console

	return Configuration{
		MySQL:                configmap.Data["mysql"],
		RedisHost:            configmap.Data["redisIP"],
		RedisPassword:        configmap.Data["redisPassword"],
		PushRoomSubscribeNum: getInt(configmap.Data, "pushRoomSubscribeNum"),
		PushAllSubscribeNum:  getInt(configmap.Data, "pushAllSubscribeNum"),

		ConnectLocalAddr:     os.Getenv("POD_IP") + RPCListenAddr,
		ConnectRPCListenAddr: RPCListenAddr,
		ConnectTCPListenAddr: ":8001",
		ConnectWSListenAddr:  ":8002",

		LogicRPCListenAddr:    RPCListenAddr,
		BusinessRPCListenAddr: RPCListenAddr,
		FileHTTPListenAddr:    "8030",

		ConnectIntClientBuilder: func() pb.ConnectIntClient {
			conn, err := grpc.DialContext(context.TODO(), k8s.GetK8STarget(namespace, "connect", RPCDialAddr), grpc.WithInsecure(),
				grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, picker.AddrPickerName)))
			if err != nil {
				panic(err)
			}
			return pb.NewConnectIntClient(conn)
		},
		LogicIntClientBuilder: func() pb.LogicIntClient {
			conn, err := grpc.DialContext(context.TODO(), k8s.GetK8STarget(namespace, "logic", RPCDialAddr), grpc.WithInsecure(),
				grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
			if err != nil {
				panic(err)
			}
			return pb.NewLogicIntClient(conn)
		},
		BusinessIntClientBuilder: func() pb.BusinessIntClient {
			conn, err := grpc.DialContext(context.TODO(), k8s.GetK8STarget(namespace, "business", RPCDialAddr), grpc.WithInsecure(),
				grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
			if err != nil {
				panic(err)
			}
			return pb.NewBusinessIntClient(conn)
		},
	}
}

func getInt(m map[string]string, key string) int {
	value, _ := strconv.Atoi(m[key])
	return value
}

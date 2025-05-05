package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"google.golang.org/grpc/balancer/roundrobin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gim/pkg/grpclib/picker"
	"gim/pkg/grpclib/resolver/k8s"
	"gim/pkg/k8sutil"
	"gim/pkg/protocol/pb/connectpb"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/protocol/pb/userpb"
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

	return Configuration{
		LogLevel: slog.LevelDebug,
		LogFile: func(server string) string {
			return fmt.Sprintf("/data/log/%s/log.log", server)
		},

		MySQL:                configmap.Data["mysql"],
		RedisHost:            configmap.Data["redisIP"],
		RedisPassword:        configmap.Data["redisPassword"],
		PushRoomSubscribeNum: getInt(configmap.Data, "pushRoomSubscribeNum"),
		PushAllSubscribeNum:  getInt(configmap.Data, "pushAllSubscribeNum"),

		ConnectLocalAddr:     os.Getenv("POD_IP") + RPCListenAddr,
		ConnectRPCListenAddr: RPCListenAddr,
		ConnectTCPListenAddr: ":8001",
		ConnectWSListenAddr:  ":8002",

		LogicRPCListenAddr: RPCListenAddr,
		UserRPCListenAddr:  RPCListenAddr,
		FileHTTPListenAddr: "8030",

		ConnectIntClientBuilder: func() connectpb.ConnectIntServiceClient {
			conn := newGrpcClient(k8s.GetK8STarget(namespace, "connect", RPCDialAddr), picker.AddrPickerName)
			return connectpb.NewConnectIntServiceClient(conn)
		},
		DeviceIntClientBuilder: func() logicpb.DeviceIntServiceClient {
			conn := newGrpcClient(k8s.GetK8STarget(namespace, "logic", RPCDialAddr), roundrobin.Name)
			return logicpb.NewDeviceIntServiceClient(conn)
		},
		MessageIntClientBuilder: func() logicpb.MessageIntServiceClient {
			conn := newGrpcClient(k8s.GetK8STarget(namespace, "logic", RPCDialAddr), roundrobin.Name)
			return logicpb.NewMessageIntServiceClient(conn)
		},
		RoomIntClientBuilder: func() logicpb.RoomIntServiceClient {
			conn := newGrpcClient(k8s.GetK8STarget(namespace, "logic", RPCDialAddr), roundrobin.Name)
			return logicpb.NewRoomIntServiceClient(conn)
		},
		UserIntClientBuilder: func() userpb.UserIntServiceClient {
			conn := newGrpcClient(k8s.GetK8STarget(namespace, "user", RPCDialAddr), roundrobin.Name)
			return userpb.NewUserIntServiceClient(conn)
		},
	}
}

func getInt(m map[string]string, key string) int {
	value, _ := strconv.Atoi(m[key])
	return value
}

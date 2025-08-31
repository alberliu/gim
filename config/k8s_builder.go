package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gim/pkg/grpclib/resolver/k8s"
	"gim/pkg/protocol/pb/businesspb"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/ugrpc"
	"gim/pkg/uk8s"
)

type k8sBuilder struct{}

func (*k8sBuilder) Build() Configuration {
	const (
		RPCListenAddr = ":8000"
		RPCDialAddr   = "8000"
	)
	const namespace = "default"

	k8sClient, err := uk8s.GetK8sClient()
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
		FileHTTPListenAddr: "8005",

		DeviceIntClientBuilder: func() logicpb.DeviceIntServiceClient {
			conn := ugrpc.NewClient(k8s.GetK8STarget(namespace, "logic", RPCDialAddr))
			return logicpb.NewDeviceIntServiceClient(conn)
		},
		MessageIntClientBuilder: func() logicpb.MessageIntServiceClient {
			conn := ugrpc.NewClient(k8s.GetK8STarget(namespace, "logic", RPCDialAddr))
			return logicpb.NewMessageIntServiceClient(conn)
		},
		RoomIntClientBuilder: func() logicpb.RoomIntServiceClient {
			conn := ugrpc.NewClient(k8s.GetK8STarget(namespace, "logic", RPCDialAddr))
			return logicpb.NewRoomIntServiceClient(conn)
		},
		UserIntClientBuilder: func() businesspb.UserIntServiceClient {
			conn := ugrpc.NewClient(k8s.GetK8STarget(namespace, "business", RPCDialAddr))
			return businesspb.NewUserIntServiceClient(conn)
		},
	}
}

func getInt(m map[string]string, key string) int {
	value, _ := strconv.Atoi(m[key])
	return value
}

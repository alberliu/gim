package config

import (
	"log/slog"

	"google.golang.org/grpc/balancer/roundrobin"

	"gim/pkg/grpclib/picker"
	_ "gim/pkg/grpclib/resolver/addrs"
	"gim/pkg/protocol/pb/connectpb"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/protocol/pb/userpb"
)

type localBuilder struct{}

func (*localBuilder) Build() Configuration {
	return Configuration{
		LogLevel: slog.LevelDebug,
		LogFile: func(server string) string {
			return ""
		},

		MySQL:                "root:123456@tcp(127.0.0.1:3306)/gim?charset=utf8mb4&parseTime=true&loc=Local",
		RedisHost:            "127.0.0.1:6379",
		RedisPassword:        "123456",
		PushRoomSubscribeNum: 100,
		PushAllSubscribeNum:  100,

		ConnectLocalAddr:     "127.0.0.1:8000",
		ConnectRPCListenAddr: ":8000",
		ConnectTCPListenAddr: ":8001",
		ConnectWSListenAddr:  ":8002",

		LogicRPCListenAddr: ":8010",
		UserRPCListenAddr:  ":8020",
		FileHTTPListenAddr: "8030",

		ConnectIntClientBuilder: func() connectpb.ConnectIntServiceClient {
			conn := newGrpcClient("addrs:///127.0.0.1:8000", picker.AddrPickerName)
			return connectpb.NewConnectIntServiceClient(conn)
		},
		DeviceIntClientBuilder: func() logicpb.DeviceIntServiceClient {
			conn := newGrpcClient("addrs:///127.0.0.1:8010", roundrobin.Name)
			return logicpb.NewDeviceIntServiceClient(conn)
		},
		MessageIntClientBuilder: func() logicpb.MessageIntServiceClient {
			conn := newGrpcClient("addrs:///127.0.0.1:8010", roundrobin.Name)
			return logicpb.NewMessageIntServiceClient(conn)
		},
		RoomIntClientBuilder: func() logicpb.RoomIntServiceClient {
			conn := newGrpcClient("addrs:///127.0.0.1:8010", roundrobin.Name)
			return logicpb.NewRoomIntServiceClient(conn)
		},
		UserIntClientBuilder: func() userpb.UserIntServiceClient {
			conn := newGrpcClient("addrs:///127.0.0.1:8020", roundrobin.Name)
			return userpb.NewUserIntServiceClient(conn)
		},
	}
}

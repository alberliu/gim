package config

import (
	"log/slog"

	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/protocol/pb/userpb"
	"gim/pkg/ugrpc"
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

		DeviceIntClientBuilder: func() logicpb.DeviceIntServiceClient {
			conn := ugrpc.NewClient("127.0.0.1:8010")
			return logicpb.NewDeviceIntServiceClient(conn)
		},
		MessageIntClientBuilder: func() logicpb.MessageIntServiceClient {
			conn := ugrpc.NewClient("127.0.0.1:8010")
			return logicpb.NewMessageIntServiceClient(conn)
		},
		RoomIntClientBuilder: func() logicpb.RoomIntServiceClient {
			conn := ugrpc.NewClient("127.0.0.1:8010")
			return logicpb.NewRoomIntServiceClient(conn)
		},
		UserIntClientBuilder: func() userpb.UserIntServiceClient {
			conn := ugrpc.NewClient("127.0.0.1:8020")
			return userpb.NewUserIntServiceClient(conn)
		},
	}
}

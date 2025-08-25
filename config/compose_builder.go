package config

import (
	"fmt"
	"log/slog"
	"net"

	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/protocol/pb/userpb"
	"gim/pkg/ugrpc"
)

type composeBuilder struct{}

func (*composeBuilder) Build() Configuration {
	addrs, err := net.LookupHost("connect")
	if err != nil {
		slog.Error("composeBuilder Build error", "error", err)
		panic(err)
	}
	if len(addrs) == 0 {
		slog.Error("composeBuilder Build error addrs is nil")
		panic(err)
	}
	connectLocalIP := addrs[0]

	return Configuration{
		LogLevel: slog.LevelDebug,
		LogFile: func(server string) string {
			return fmt.Sprintf("/data/log/%s/log.log", server)
		},

		MySQL:                "root:123456@tcp(mysql:3306)/gim?charset=utf8mb4&parseTime=true&loc=Local",
		RedisHost:            "redis:6379",
		RedisPassword:        "123456",
		PushRoomSubscribeNum: 100,
		PushAllSubscribeNum:  100,

		ConnectLocalAddr:     connectLocalIP + ":8000",
		ConnectRPCListenAddr: ":8000",
		ConnectTCPListenAddr: ":8001",
		ConnectWSListenAddr:  ":8002",

		LogicRPCListenAddr: ":8010",
		UserRPCListenAddr:  ":8020",
		FileHTTPListenAddr: "8030",

		DeviceIntClientBuilder: func() logicpb.DeviceIntServiceClient {
			conn := ugrpc.NewClient("dns:///logic:8010")
			return logicpb.NewDeviceIntServiceClient(conn)
		},
		MessageIntClientBuilder: func() logicpb.MessageIntServiceClient {
			conn := ugrpc.NewClient("dns:///logic:8010")
			return logicpb.NewMessageIntServiceClient(conn)
		},
		RoomIntClientBuilder: func() logicpb.RoomIntServiceClient {
			conn := ugrpc.NewClient("dns:///logic:8010")
			return logicpb.NewRoomIntServiceClient(conn)
		},

		UserIntClientBuilder: func() userpb.UserIntServiceClient {
			conn := ugrpc.NewClient("dns:///user:8020")
			return userpb.NewUserIntServiceClient(conn)
		},
	}
}

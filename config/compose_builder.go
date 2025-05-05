package config

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc/balancer/roundrobin"

	"gim/pkg/grpclib/picker"
	_ "gim/pkg/grpclib/resolver/addrs"
	"gim/pkg/protocol/pb/connectpb"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/protocol/pb/userpb"
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

		ConnectIntClientBuilder: func() connectpb.ConnectIntServiceClient {
			conn := newGrpcClient("dns:///connect:8000", picker.AddrPickerName)
			return connectpb.NewConnectIntServiceClient(conn)
		},

		DeviceIntClientBuilder: func() logicpb.DeviceIntServiceClient {
			conn := newGrpcClient("dns:///logic:8010", roundrobin.Name)
			return logicpb.NewDeviceIntServiceClient(conn)
		},
		MessageIntClientBuilder: func() logicpb.MessageIntServiceClient {
			conn := newGrpcClient("dns:///logic:8010", roundrobin.Name)
			return logicpb.NewMessageIntServiceClient(conn)
		},
		RoomIntClientBuilder: func() logicpb.RoomIntServiceClient {
			conn := newGrpcClient("dns:///logic:8010", roundrobin.Name)
			return logicpb.NewRoomIntServiceClient(conn)
		},

		UserIntClientBuilder: func() userpb.UserIntServiceClient {
			conn := newGrpcClient("dns:///user:8020", roundrobin.Name)
			return userpb.NewUserIntServiceClient(conn)
		},
	}
}

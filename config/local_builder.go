package config

import (
	"log/slog"
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

		LogicRPCListenAddr:    ":8010",
		BusinessRPCListenAddr: ":8020",
		FileHTTPListenAddr:    "8030",

		LogicServerAddr:    "127.0.0.1:8010",
		BusinessServerAddr: "127.0.0.1:8020",
	}
}

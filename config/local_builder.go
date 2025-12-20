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
	}
}

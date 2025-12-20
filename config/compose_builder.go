package config

import (
	"fmt"
	"log/slog"
)

type composeBuilder struct{}

func (*composeBuilder) Build() Configuration {
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
	}
}

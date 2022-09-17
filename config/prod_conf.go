package config

import (
	"gim/pkg/logger"

	"go.uber.org/zap"
)

func initProdConf() {
	Config = config{
		MySQL:         "root:gim123456@tcp(111.229.238.28:3306)/gim?charset=utf8&parseTime=true",
		RedisIP:       "111.229.238.28:6379",
		RedisPassword: "alber123456",

		TCPListenAddr:        ":8080",
		WSListenAddr:         ":8081",
		ConnectRPCListenAddr: ":50000",
		ConnectLocalAddr:     "127.0.0.1:50000",
		PushRoomSubscribeNum: 100,
		PushAllSubscribeNum:  100,

		LogicRPCListenAddr: ":50100",

		BusinessRPCListenAddr: ":50200",
	}

	logger.Level = zap.DebugLevel
	logger.Target = logger.Console
}

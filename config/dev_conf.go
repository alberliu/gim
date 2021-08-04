package config

import (
	"gim/pkg/logger"

	"go.uber.org/zap"
)

func initDevConf() {
	RPCAddr = RPCAddrConf{
		ConnectRPCAddr:  "addrs:///127.0.0.1:50000",
		LogicRPCAddr:    "addrs:///127.0.0.1:50100",
		BusinessRPCAddr: "addrs:///127.0.0.1:50200",
	}

	Connect = ConnectConf{
		TCPListenAddr: ":8080",
		WSListenAddr:  ":8081",
		RPCListenAddr: ":50000",
		LocalAddr:     "127.0.0.1:50000",
		RedisIP:       "111.229.238.28:6379",
		RedisPassword: "alber123456",
		SubscribeNum:  100,
	}

	Logic = LogicConf{
		MySQL:         "root:gim123456@tcp(111.229.238.28:3306)/gim?charset=utf8&parseTime=true",
		NSQIP:         "111.229.238.28:4150",
		RedisIP:       "111.229.238.28:6379",
		RedisPassword: "alber123456",
		RPCListenAddr: ":50100",
	}

	Business = BusinessConf{
		MySQL:         "root:gim123456@tcp(111.229.238.28:3306)/gim?charset=utf8&parseTime=true",
		NSQIP:         "111.229.238.28:4150",
		RedisIP:       "111.229.238.28:6379",
		RedisPassword: "alber123456",
		RPCListenAddr: ":50200",
	}

	logger.Leavel = zap.DebugLevel
	logger.Target = logger.File
}

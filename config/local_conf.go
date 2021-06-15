package config

import (
	"gim/pkg/logger"

	"go.uber.org/zap"
)

func initLocalConf() {
	Logic = LogicConf{
		MySQL:            "root:gim123456@tcp(111.229.238.28:3306)/gim?charset=utf8&parseTime=true",
		NSQIP:            "111.229.238.28:4150",
		RedisIP:          "111.229.238.28:6379",
		RedisPassword:    "alber123456",
		RPCIntListenAddr: ":50000",
		RPCExtListenAddr: ":50001",
		ConnectRPCAddrs:  "addrs:///127.0.0.1:50100,127.0.0.1:50200",
		BusinessRPCAddrs: "addrs:///127.0.0.1:50300",
	}

	Connect = ConnectConf{
		TCPListenAddr: 8080,
		WSListenAddr:  ":8081",
		RPCListenAddr: ":50100",
		LocalAddr:     "127.0.0.1:50100",
		LogicRPCAddrs: "addrs:///127.0.0.1:50000",
		RedisIP:       "111.229.238.28:6379",
		RedisPassword: "alber123456",
		SubscribeNum:  100,
	}

	Business = BusinessConf{
		MySQL:            "root:gim123456@tcp(111.229.238.28:3306)/gim?charset=utf8&parseTime=true",
		NSQIP:            "111.229.238.28:4150",
		RedisIP:          "111.229.238.28:6379",
		RedisPassword:    "alber123456",
		RPCIntListenAddr: ":50300",
		RPCExtListenAddr: ":50301",
		LogicRPCAddrs:    "addrs:///127.0.0.1:50000",
	}
	logger.Leavel = zap.DebugLevel
	logger.Target = logger.Console
}

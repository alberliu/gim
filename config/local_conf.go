package config

import (
	"gim/pkg/logger"

	"go.uber.org/zap"
)

func initLocalConf() {
	Logic = LogicConf{
		MySQL:            "root:liu123456@tcp(localhost:3307)/gim?charset=utf8&parseTime=true&loc=Local",
		NSQIP:            "127.0.0.1:4150",
		RedisIP:          "127.0.0.1:6379",
		RedisPassword:    "liu123456",
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
		RedisPassword: "",
		SubscribeNum:  100,
	}

	Business = BusinessConf{
		MySQL:            "root:liu123456@tcp(localhost:3306)/im?charset=utf8&parseTime=true",
		NSQIP:            "127.0.0.1:4150",
		RedisIP:          "127.0.0.1:6379",
		RPCIntListenAddr: ":50300",
		RPCExtListenAddr: ":50301",
		LogicRPCAddrs:    "addrs:///127.0.0.1:50000",
	}

	logger.Leavel = zap.DebugLevel
	logger.Target = logger.Console
}

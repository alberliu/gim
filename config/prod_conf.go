package config

import (
	"gim/pkg/logger"

	"go.uber.org/zap"
)

func initProdConf() {
	Logic = LogicConf{
		MySQL:            "root:liu123456@tcp(localhost:3306)/gim?charset=utf8&parseTime=true&loc=Local",
		NSQIP:            "127.0.0.1:4150",
		RedisIP:          "127.0.0.1:6379",
		RedisPassword:    "liu123456",
		RPCIntListenAddr: ":50000",
		RPCExtListenAddr: ":50001",
		ConnRPCAddrs:     "addrs:///127.0.0.1:50100,127.0.0.1:50200",
		UserRPCAddrs:     "addrs:///127.0.0.1:50300",
	}

	TCPConn = TCPConnConf{
		TCPListenAddr: ":8080",
		RPCListenAddr: ":50100",
		LocalAddr:     "127.0.0.1:50100",
		LogicRPCAddrs: "addrs:///127.0.0.1:50000",
	}

	WSConn = WSConnConf{
		WSListenAddr:  ":8081",
		RPCListenAddr: ":50200",
		LocalAddr:     "127.0.0.1:50200",
		LogicRPCAddrs: "addrs:///127.0.0.1:50000",
	}

	User = UserConf{
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

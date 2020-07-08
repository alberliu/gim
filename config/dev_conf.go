package config

import (
	"gim/pkg/logger"

	"go.uber.org/zap"
)

func initDevConf() {
	LogicConf = logicConf{
		MySQL:                  "root:liu123456@tcp(localhost:3306)/gim?charset=utf8&parseTime=true",
		NSQIP:                  "127.0.0.1:4150",
		RedisIP:                "127.0.0.1:6379",
		RPCIntListenAddr:       ":50000",
		ClientRPCExtListenAddr: ":50001",
		ServerRPCExtListenAddr: ":50002",
		ConnRPCAddrs:           "addrs:///127.0.0.1:60000,127.0.0.1:60001",
	}

	TCPConnConf = tcpConnConf{
		Port:          8080,
		RPCListenAddr: ":60000",
		LocalAddr:     "127.0.0.1:60000",
		LogicRPCAddrs: "addrs:///127.0.0.1:50000",
	}

	WSConnConf = wsConnConf{
		WSListenAddr:  ":8081",
		RPCListenAddr: ":60001",
		LocalAddr:     "127.0.0.1:60001",
		LogicRPCAddrs: "addrs:///127.0.0.1:50000",
	}

	logger.Leavel = zap.DebugLevel
	logger.Target = logger.File
}

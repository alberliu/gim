package config

import (
	"gim/pkg/logger"

	"go.uber.org/zap"
)

func initDevConf() {
	Logic = LogicConf{
		MySQL:            "root:gim123456@tcp(111.229.238.28:3306)/gim?charset=utf8&parseTime=true",
		NSQIP:            "111.229.238.28:4150",
		RedisIP:          "111.229.238.28:6379",
		RedisPassword:    "",
		RPCIntListenAddr: ":50000",
		RPCExtListenAddr: ":50001",
		ConnRPCAddrs:     "addrs:///127.0.0.1:50100,127.0.0.1:50200",
		BusinessRPCAddrs: "addrs:///127.0.0.1:50300",
	}
	TCPConn = TCPConnConf{
		TCPListenAddr: 8080,
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
	Business = BusinessConf{
		MySQL:            "root:gim123456@tcp(111.229.238.28:3306)/gim?charset=utf8&parseTime=true",
		NSQIP:            "111.229.238.28:4150",
		RedisIP:          "111.229.238.28:6379",
		RedisPassword:    "",
		RPCIntListenAddr: ":50300",
		RPCExtListenAddr: ":50301",
		LogicRPCAddrs:    "addrs:///127.0.0.1:50000",
	}
	logger.Leavel = zap.DebugLevel
	logger.Target = logger.File
}

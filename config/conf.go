package config

import (
	"os"
)

var (
	Logic   LogicConf
	TCPConn TCPConnConf
	WSConn  WSConnConf
	User    UserConf
)

// logic配置
type LogicConf struct {
	MySQL            string
	NSQIP            string
	RedisIP          string
	RedisPassword    string
	RPCIntListenAddr string
	RPCExtListenAddr string
	ConnRPCAddrs     string
	UserRPCAddrs     string
}

// TCPConnConf配置
type TCPConnConf struct {
	TCPListenAddr int
	RPCListenAddr string
	LocalAddr     string
	LogicRPCAddrs string
}

// WS配置
type WSConnConf struct {
	WSListenAddr  string
	RPCListenAddr string
	LocalAddr     string
	LogicRPCAddrs string
}

// User配置
type UserConf struct {
	MySQL            string
	NSQIP            string
	RedisIP          string
	RedisPassword    string
	RPCIntListenAddr string
	RPCExtListenAddr string
	LogicRPCAddrs    string
}

func init() {
	env := os.Getenv("im_env")
	switch env {
	case "dev":
		initDevConf()
	case "prod":
		initProdConf()
	default:
		initLocalConf()
	}
}

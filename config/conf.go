package config

import (
	"os"
)

var (
	LogicConf   logicConf
	TCPConnConf tcpConnConf
	WSConnConf  wsConnConf
)

// logic配置
type logicConf struct {
	MySQL                  string
	NSQIP                  string
	RedisIP                string
	RPCIntListenAddr       string
	ClientRPCExtListenAddr string
	ServerRPCExtListenAddr string
	ConnRPCAddrs           string
}

// conn配置
type tcpConnConf struct {
	Port          int
	RPCListenAddr string
	LocalAddr     string
	LogicRPCAddrs string
}

// WS配置
type wsConnConf struct {
	WSListenAddr  string
	RPCListenAddr string
	LocalAddr     string
	LogicRPCAddrs string
}

func init() {
	env := os.Getenv("gim_env")
	switch env {
	case "dev":
		initDevConf()
	case "pre":
		initPreConf()
	case "prod":
		initProdConf()
	default:
		initLocalConf()
	}
}

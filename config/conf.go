package config

import (
	"os"
)

var (
	LogicConf logicConf
	ConnConf  connConf
	WSConf    wsConf
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
type connConf struct {
	TCPListenAddr string
	RPCListenAddr string
	LocalAddr     string
	LogicRPCAddrs string
}

// WS配置
type wsConf struct {
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

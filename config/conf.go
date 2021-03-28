package config

import (
	"os"
)

var (
	Logic    LogicConf
	Connect  ConnectConf
	Business BusinessConf
)

// logic配置
type LogicConf struct {
	MySQL            string
	NSQIP            string
	RedisIP          string
	RedisPassword    string
	RPCIntListenAddr string
	RPCExtListenAddr string
	ConnectRPCAddrs  string
	BusinessRPCAddrs string
}

// TCPConnConf配置
type ConnectConf struct {
	TCPListenAddr int
	WSListenAddr  string
	RPCListenAddr string
	LocalAddr     string
	LogicRPCAddrs string
	RedisIP       string
	RedisPassword string
	SubscribeNum  int
}

// Business配置
type BusinessConf struct {
	MySQL            string
	NSQIP            string
	RedisIP          string
	RedisPassword    string
	RPCIntListenAddr string
	RPCExtListenAddr string
	LogicRPCAddrs    string
}

func init() {
	env := os.Getenv("gim_env")
	switch env {
	case "dev":
		initDevConf()
	case "prod":
		initProdConf()
	default:
		initLocalConf()
	}
}

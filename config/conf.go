package config

import (
	"os"
)

var (
	RPCAddr  RPCAddrConf
	Logic    LogicConf
	Connect  ConnectConf
	Business BusinessConf
)

// RPCAddrConf RPC配置
type RPCAddrConf struct {
	ConnectRPCAddr  string
	BusinessRPCAddr string
	LogicRPCAddr    string
}

// ConnectConf Connect配置
type ConnectConf struct {
	TCPListenAddr string
	WSListenAddr  string
	RPCListenAddr string
	LocalAddr     string
	RedisIP       string
	RedisPassword string
	SubscribeNum  int
}

// LogicConf logic配置
type LogicConf struct {
	MySQL         string
	NSQIP         string
	RedisIP       string
	RedisPassword string
	RPCListenAddr string
}

// BusinessConf Business配置
type BusinessConf struct {
	MySQL         string
	NSQIP         string
	RedisIP       string
	RedisPassword string
	RPCListenAddr string
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

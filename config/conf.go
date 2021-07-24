package config

import (
	"os"
)

var (
	Logic    LogicConf
	Connect  ConnectConf
	Business BusinessConf
)

// LogicConf logic配置
type LogicConf struct {
	MySQL            string
	NSQIP            string
	RedisIP          string
	RedisPassword    string
	RPCListenAddr    string
	ConnectRPCAddrs  string
	BusinessRPCAddrs string
}

// ConnectConf Connect配置
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

// BusinessConf Business配置
type BusinessConf struct {
	MySQL         string
	NSQIP         string
	RedisIP       string
	RedisPassword string
	RPCListenAddr string
	LogicRPCAddrs string
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

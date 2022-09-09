package config

import (
	"os"
)

var Config config

type config struct {
	MySQL         string
	RedisIP       string
	RedisPassword string

	ConnectRPCAddr  string
	BusinessRPCAddr string
	LogicRPCAddr    string

	TCPListenAddr        string
	WSListenAddr         string
	ConnectRPCListenAddr string
	ConnectLocalAddr     string
	PushRoomSubscribeNum int
	PushAllSubscribeNum  int

	LogicRPCListenAddr    string
	BusinessRPCListenAddr string
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

package main

import (
	"gim/conf"
	"gim/conn"
	"gim/public/rpc_cli"
	"gim/public/util"
)

func main() {
	// 启动rpc服务
	go func() {
		defer util.RecoverPanic()
		conn.StartRPCServer()
	}()

	// 初始化Rpc Client
	rpc_cli.InitLogicIntClient()

	// 启动长链接服务器
	server := conn.NewTCPServer(conf.ConnTCPListenAddr, 1)
	server.Start()
}

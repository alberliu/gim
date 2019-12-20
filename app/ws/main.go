package main

import (
	"gim/conf"
	"gim/public/rpc_cli"
	"gim/public/util"
	"gim/ws"
)

func main() {
	// 启动rpc服务
	go func() {
		defer util.RecoverPanic()
		ws.StartRPCServer()
	}()

	// 初始化Rpc Client
	rpc_cli.InitLogicIntClient()

	// 启动长链接服务器
	ws.StartWSServer(conf.ConnTCPListenAddr)
}

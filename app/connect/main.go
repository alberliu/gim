package main

import (
	"gim/conf"
	"gim/connect"
	"gim/public/util"
)

func main() {
	// 启动rpc服务
	go func() {
		defer util.RecoverPanic()
		connect.StartRPCServer()
	}()

	// 初始化Rpc Client
	go func() {
		defer util.RecoverPanic()
		connect.InitRpcClient()
	}()

	// 启动长链接服务器
	server := connect.NewTCPServer(conf.ConnectTCPListenAddress, 1)
	server.Start()
}

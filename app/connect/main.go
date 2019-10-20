package main

import (
	"goim/conf"
	"goim/connect"
	"goim/public/util"
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
	conf := connect.Conf{
		Address:      conf.ConnectTCPListenIP + ":" + conf.ConnectTCPListenPort,
		MaxConnCount: 100,
		AcceptCount:  1,
	}
	server := connect.NewTCPServer(conf)
	server.Start()
}

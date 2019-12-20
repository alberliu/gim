package main

import (
	"gim/conf"
	"gim/public/rpc_cli"
	"gim/public/util"
	"gim/tcp_conn"
)

func main() {
	// 启动rpc服务
	go func() {
		defer util.RecoverPanic()
		tcp_conn.StartRPCServer()
	}()

	// 初始化Rpc Client
	rpc_cli.InitLogicIntClient(conf.ConnConf.LogicRPCAddrs)

	// 启动长链接服务器
	server := tcp_conn.NewTCPServer(conf.ConnConf.TCPListenAddr, 10)
	server.Start()
}

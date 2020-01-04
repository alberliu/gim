package main

import (
	"gim/api/tcp_conn"
	"gim/config"
	tcp_conn2 "gim/internal/tcp_conn"
	"gim/pkg/rpc_cli"
	"gim/pkg/util"
)

func main() {
	// 启动rpc服务
	go func() {
		defer util.RecoverPanic()
		tcp_conn.StartRPCServer()
	}()

	// 初始化Rpc Client
	rpc_cli.InitLogicIntClient(config.ConnConf.LogicRPCAddrs)

	// 启动长链接服务器
	server := tcp_conn2.NewTCPServer(config.ConnConf.TCPListenAddr, 10)
	server.Start()
}

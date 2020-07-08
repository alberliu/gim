package main

import (
	"gim/config"
	"gim/internal/tcp_conn"
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
	rpc_cli.InitLogicIntClient(config.TCPConnConf.LogicRPCAddrs)

	// 启动长链接服务器
	// 启动长链接服务器
	tcp_conn.StartTCPServer(config.TCPConnConf.Port)
}

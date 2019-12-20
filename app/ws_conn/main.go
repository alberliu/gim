package main

import (
	"gim/conf"
	"gim/public/rpc_cli"
	"gim/public/util"
	"gim/ws_conn"
)

func main() {
	// 启动rpc服务
	go func() {
		defer util.RecoverPanic()
		ws_conn.StartRPCServer()
	}()

	// 初始化Rpc Client
	rpc_cli.InitLogicIntClient(conf.WSConf.LogicRPCAddrs)

	// 启动长链接服务器
	ws_conn.StartWSServer(conf.WSConf.WSListenAddr)
}

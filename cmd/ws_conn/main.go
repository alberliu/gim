package main

import (
	"gim/config"
	"gim/internal/ws_conn"
	"gim/pkg/logger"
	"gim/pkg/rpc"
	"gim/pkg/util"
)

func main() {
	logger.Init()
	// 启动rpc服务
	go func() {
		defer util.RecoverPanic()
		ws_conn.StartRPCServer()
	}()

	// 初始化Rpc Client
	rpc.InitLogicIntClient(config.WSConn.LogicRPCAddrs)

	// 启动长链接服务器
	ws_conn.StartWSServer(config.WSConn.WSListenAddr)
}

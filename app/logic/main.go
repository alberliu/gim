package main

import (
	"gim/logic/db"
	"gim/logic/server"
	"gim/public/logger"
	"gim/public/rpc_cli"
	"gim/public/util"
)

func main() {
	// 初始化数据库
	db.InitDB()

	// 初始化自增id配置
	util.InitUID(db.DBCli)

	// 初始化RpcClient
	rpc_cli.InitConnIntClient()

	/*// 启动nsq消费服务
	go func() {
		defer util.RecoverPanic()
		consume.StartNsqConsumer()
	}()
	*/

	server.StartRpcServer()
	logger.Logger.Info("logic server start")
	select {}
}

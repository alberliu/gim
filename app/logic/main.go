package main

import (
	"gim/conf"
	"gim/logic/db"
	"gim/logic/rpc"
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
	rpc_cli.InitConnIntClient(conf.LogicConf.ConnRPCAddrs)

	/*// 启动nsq消费服务
	go func() {
		defer util.RecoverPanic()
		consume.StartNsqConsumer()
	}()
	*/

	rpc.StartRpcServer()
	logger.Logger.Info("logic server start")
	select {}
}

package main

import (
	"gim/conf"
	"gim/logic/controller"
	"gim/logic/db"
	"gim/logic/rpc/client"
	"gim/logic/rpc/server"
	"gim/public/logger"
	"gim/public/util"
)

func main() {
	// 初始化数据库
	db.InitDB()

	// 初始化自增id配置
	util.InitUID(db.DBCli)

	// 启动rpc服务
	go func() {
		defer util.RecoverPanic()
		server.StartRPCServer()
	}()

	// 初始化RpcClient
	go func() {
		defer util.RecoverPanic()
		client.InitRpcClient()
	}()

	/*// 启动nsq消费服务
	go func() {
		defer util.RecoverPanic()
		consume.StartNsqConsumer()
	}()
	*/

	// 启动web容器
	err := controller.Engine.Run(conf.LogicHTTPListenIP)
	if err != nil {
		logger.Sugar.Error(err)
	}
}

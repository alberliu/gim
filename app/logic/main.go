package main

import (
	"goim/conf"
	"goim/logic/controller"
	"goim/logic/db"
	"goim/logic/rpc/client"
	"goim/logic/rpc/server"
	"goim/public/logger"
	"goim/public/util"
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

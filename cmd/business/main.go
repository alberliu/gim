package main

import (
	"gim/config"
	"gim/internal/business/api"
	"gim/pkg/db"
	"gim/pkg/logger"
	"gim/pkg/rpc"
)

func main() {
	logger.Init()
	db.InitMysql(config.Business.MySQL)
	db.InitRedis(config.Business.RedisIP, config.Logic.RedisPassword)

	// 初始化RpcClient
	rpc.InitLogicIntClient(config.Business.LogicRPCAddrs)

	api.StartRpcServer()
	logger.Logger.Info("user server start")
	select {}
}

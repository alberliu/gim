package main

import (
	"gim/config"
	"gim/internal/user/api"
	"gim/pkg/db"
	"gim/pkg/logger"
	"gim/pkg/rpc"
)

func main() {
	logger.Init()
	db.InitMysql(config.User.MySQL)
	db.InitRedis(config.User.RedisIP, config.Logic.RedisPassword)

	// 初始化RpcClient
	rpc.InitLogicIntClient(config.User.LogicRPCAddrs)

	api.StartRpcServer()
	logger.Logger.Info("user server start")
	select {}
}

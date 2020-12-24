package main

import (
	"gim/config"
	"gim/internal/logic/api"
	"gim/pkg/db"
	"gim/pkg/logger"
	"gim/pkg/rpc"
)

func main() {
	logger.Init()
	db.InitMysql(config.Logic.MySQL)
	db.InitRedis(config.Logic.RedisIP, config.Logic.RedisPassword)

	// 初始化RpcClient
	rpc.InitConnIntClient(config.Logic.ConnRPCAddrs)
	rpc.InitBusinessIntClient(config.Logic.UserRPCAddrs)

	api.StartRpcServer()
	logger.Logger.Info("logic server start")
	select {}
}

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
	rpc.InitConnectIntClient(config.Logic.ConnectRPCAddrs)
	rpc.InitBusinessIntClient(config.Logic.BusinessRPCAddrs)

	api.StartRpcServer()
	logger.Logger.Info("logic server start")
	select {}
}

package main

import (
	"gim/config"
	"gim/internal/logic/api"
	"gim/internal/logic/db"
	"gim/pkg/logger"
	"gim/pkg/rpc_cli"
	"gim/pkg/util"
)

func main() {
	// 初始化数据库
	db.InitDB()

	// 初始化自增id配置
	util.InitUID(db.DBCli)

	// 初始化RpcClient
	rpc_cli.InitConnIntClient(config.LogicConf.ConnRPCAddrs)

	api.StartRpcServer()
	logger.Logger.Info("logic server start")
	select {}
}

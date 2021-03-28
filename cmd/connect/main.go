package main

import (
	"context"
	"gim/config"
	"gim/internal/connect"
	"gim/pkg/db"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc"
	"gim/pkg/util"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	logger.Init()

	db.InitRedis(config.Connect.RedisIP, config.Connect.RedisPassword)
	// 启动rpc服务
	go func() {
		defer util.RecoverPanic()
		connect.StartRPCServer()
	}()

	// 初始化Rpc Client
	rpc.InitLogicIntClient(config.Connect.LogicRPCAddrs)

	// 启动TCP长链接服务器
	go func() {
		connect.StartTCPServer()
	}()

	// 启动WebSocket长链接服务器
	go func() {
		defer util.RecoverPanic()
		connect.StartWSServer(config.Connect.WSListenAddr)
	}()

	// 启动服务订阅
	connect.StartSubscribe()

	c := make(chan os.Signal, 0)
	signal.Notify(c, syscall.SIGTERM)

	s := <-c
	logger.Logger.Info("server stop start", zap.Any("signal", s))
	rpc.LogicIntClient.ServerStop(context.TODO(), &pb.ServerStopReq{ConnAddr: config.Connect.LocalAddr})
	logger.Logger.Info("server stop end")
}

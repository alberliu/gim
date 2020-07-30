package main

import (
	"context"
	"gim/config"
	"gim/internal/ws_conn"
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
	// 启动rpc服务
	go func() {
		defer util.RecoverPanic()
		ws_conn.StartRPCServer()
	}()

	// 初始化Rpc Client
	rpc.InitLogicIntClient(config.WSConn.LogicRPCAddrs)

	// 启动长链接服务器
	go func() {
		defer util.RecoverPanic()
		ws_conn.StartWSServer(config.WSConn.WSListenAddr)
	}()

	c := make(chan os.Signal, 0)
	signal.Notify(c, syscall.SIGTERM)

	s := <-c
	logger.Logger.Info("server stop start", zap.Any("signal", s))
	rpc.LogicIntClient.ServerStop(context.TODO(), &pb.ServerStopReq{ConnAddr: config.WSConn.LocalAddr})
	logger.Logger.Info("server stop end")
}

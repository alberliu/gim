package main

import (
	"context"
	"gim/config"
	"gim/internal/tcp_conn"
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
		tcp_conn.StartRPCServer()
	}()

	// 初始化Rpc Client
	rpc.InitLogicIntClient(config.TCPConn.LogicRPCAddrs)

	// 启动长链接服务器
	go func() {
		tcp_conn.StartTCPServer()
	}()

	c := make(chan os.Signal, 0)
	signal.Notify(c, syscall.SIGTERM)

	s := <-c
	logger.Logger.Info("server stop start", zap.Any("signal", s))
	rpc.LogicIntClient.ServerStop(context.TODO(), &pb.ServerStopReq{ConnAddr: config.TCPConn.LocalAddr})
	logger.Logger.Info("server stop end")
}

package main

import (
	"context"
	"gim/config"
	"gim/internal/connect"
	"gim/pkg/db"
	"gim/pkg/interceptor"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"go.uber.org/zap"
)

func main() {
	logger.Init()

	db.InitRedis(config.Connect.RedisIP, config.Connect.RedisPassword)

	// 初始化Rpc Client
	rpc.InitLogicIntClient(config.RPCAddr.LogicRPCAddr)

	// 启动TCP长链接服务器
	go func() {
		connect.StartTCPServer()
	}()

	// 启动WebSocket长链接服务器
	go func() {
		connect.StartWSServer(config.Connect.WSListenAddr)
	}()

	// 启动服务订阅
	connect.StartSubscribe()

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.NewInterceptor("connect_interceptor", nil)))

	// 监听服务关闭信号，服务平滑重启
	go func() {
		c := make(chan os.Signal, 0)
		signal.Notify(c, syscall.SIGTERM)
		s := <-c
		logger.Logger.Info("server stop start", zap.Any("signal", s))
		rpc.LogicIntClient.ServerStop(context.TODO(), &pb.ServerStopReq{ConnAddr: config.Connect.LocalAddr})
		logger.Logger.Info("server stop end")

		server.GracefulStop()
	}()

	pb.RegisterConnectIntServer(server, &connect.ConnIntServer{})
	listener, err := net.Listen("tcp", config.Connect.RPCListenAddr)
	if err != nil {
		panic(err)
	}

	logger.Logger.Info("rpc服务已经开启")
	err = server.Serve(listener)
	if err != nil {
		logger.Logger.Error("Serve", zap.Error(err))
	}
}

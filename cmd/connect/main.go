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
	config.Init()
	db.Init()

	// 启动TCP长链接服务器
	go func() {
		connect.StartTCPServer(config.TCPListenAddr)
	}()

	// 启动WebSocket长链接服务器
	go func() {
		connect.StartWSServer(config.WSListenAddr)
	}()

	// 启动服务订阅
	connect.StartSubscribe()

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.NewInterceptor("connect_interceptor", nil)))

	// 监听服务关闭信号，服务平滑重启
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		s := <-c
		logger.Logger.Info("server stop start", zap.Any("signal", s))
		_, _ = rpc.GetLogicIntClient().ServerStop(context.TODO(), &pb.ServerStopReq{ConnAddr: config.LocalAddr})
		logger.Logger.Info("server stop end")

		server.GracefulStop()
	}()

	pb.RegisterConnectIntServer(server, &connect.ConnIntServer{})
	listener, err := net.Listen("tcp", config.RPCListenAddr)
	if err != nil {
		panic(err)
	}

	logger.Logger.Info("rpc服务已经开启")
	err = server.Serve(listener)
	if err != nil {
		logger.Logger.Error("serve error", zap.Error(err))
	}
}

package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"gim/config"
	"gim/internal/connect"
	"gim/pkg/interceptor"
	"gim/pkg/logger"
	"gim/pkg/protocol/pb"
	"gim/pkg/rpc"
)

func main() {
	// 启动TCP长链接服务器
	go func() {
		connect.StartTCPServer(config.Config.ConnectTCPListenAddr)
	}()

	// 启动WebSocket长链接服务器
	go func() {
		connect.StartWSServer(config.Config.ConnectWSListenAddr)
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
		_, _ = rpc.GetLogicIntClient().ServerStop(context.TODO(), &pb.ServerStopReq{ConnAddr: config.Config.ConnectLocalAddr})
		logger.Logger.Info("server stop end")

		server.GracefulStop()
	}()

	pb.RegisterConnectIntServer(server, &connect.ConnIntServer{})
	listener, err := net.Listen("tcp", config.Config.ConnectRPCListenAddr)
	if err != nil {
		panic(err)
	}

	logger.Logger.Info("rpc服务已经开启")
	err = server.Serve(listener)
	if err != nil {
		logger.Logger.Error("serve error", zap.Error(err))
	}
}

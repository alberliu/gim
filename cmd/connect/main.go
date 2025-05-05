package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"gim/config"
	"gim/internal/connect"
	"gim/pkg/interceptor"
	"gim/pkg/logger"
	pb "gim/pkg/protocol/pb/connectpb"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/rpc"
)

func main() {
	logger.Init("connect")

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

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.NewInterceptor(nil)))

	// 监听服务关闭信号，服务平滑重启
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		s := <-c
		slog.Info("server stop start", "signal", s)
		_, _ = rpc.GetDeviceIntClient().ServerStop(context.TODO(), &logicpb.ServerStopRequest{ConnAddr: config.Config.ConnectLocalAddr})
		slog.Info("server stop end")

		server.GracefulStop()
	}()

	pb.RegisterConnectIntServiceServer(server, &connect.ConnIntService{})
	listener, err := net.Listen("tcp", config.Config.ConnectRPCListenAddr)
	if err != nil {
		panic(err)
	}

	slog.Info("rpc服务已经开启")
	err = server.Serve(listener)
	if err != nil {
		slog.Error("serve error", "error", err)
	}
}

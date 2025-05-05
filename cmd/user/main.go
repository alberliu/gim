package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"gim/config"
	"gim/internal/user/api"
	"gim/pkg/interceptor"
	"gim/pkg/logger"
	pb "gim/pkg/protocol/pb/userpb"
	"gim/pkg/urlwhitelist"
)

func main() {
	logger.Init("user")

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.NewInterceptor(urlwhitelist.User)))

	// 监听服务关闭信号，服务平滑重启
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		s := <-c
		slog.Info("server stop", "signal", s)
		server.GracefulStop()
	}()

	pb.RegisterUserIntServiceServer(server, &api.UserIntService{})
	pb.RegisterUserExtServiceServer(server, &api.UserExtService{})
	listen, err := net.Listen("tcp", config.Config.UserRPCListenAddr)
	if err != nil {
		panic(err)
	}

	slog.Info("rpc服务已经开启")
	err = server.Serve(listen)
	if err != nil {
		slog.Error("serve error", "error", err)
	}
}

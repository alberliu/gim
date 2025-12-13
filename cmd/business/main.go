package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"gim/config"
	friendapi "gim/internal/business/friend/api"
	userapi "gim/internal/business/user/api"
	"gim/pkg/interceptor"
	"gim/pkg/logger"
	pb "gim/pkg/protocol/pb/businesspb"
)

func main() {
	logger.Init("user")

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.NewInterceptor(interceptor.UserWhitelistURL)))

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	pb.RegisterUserIntServiceServer(server, &userapi.UserIntService{})
	pb.RegisterUserExtServiceServer(server, &userapi.UserExtService{})
	pb.RegisterFriendExtServiceServer(server, &friendapi.FriendExtService{})
	listen, err := net.Listen("tcp", config.Config.BusinessRPCListenAddr)
	if err != nil {
		panic(err)
	}

	go func() {
		slog.Info("rpc服务已经开启")
		err = server.Serve(listen)
		if err != nil {
			slog.Error("serve error", "error", err)
			panic(err)
		}
	}()

	// 监听服务关闭信号，服务平滑重启
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	s := <-c
	slog.Info("server stop", "signal", s)
	server.GracefulStop()
}

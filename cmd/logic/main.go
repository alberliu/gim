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
	deviceapi "gim/internal/logic/device/api"
	groupapi "gim/internal/logic/group/api"
	messageapi "gim/internal/logic/message/api"
	"gim/internal/logic/room"
	"gim/pkg/interceptor"
	"gim/pkg/logger"
	pb "gim/pkg/protocol/pb/logicpb"
)

func main() {
	logger.Init("logic")

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.NewInterceptor(nil)))

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	pb.RegisterDeviceIntServiceServer(server, &deviceapi.DeviceIntService{})
	pb.RegisterMessageExtServiceServer(server, &messageapi.MessageExtService{})
	pb.RegisterMessageIntServiceServer(server, &messageapi.MessageIntService{})
	pb.RegisterGroupIntServiceServer(server, &groupapi.GroupIntService{})
	pb.RegisterRoomIntServiceServer(server, &room.RoomIntService{})

	listen, err := net.Listen("tcp", config.Config.LogicRPCListenAddr)
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

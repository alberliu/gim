package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

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

	// 监听服务关闭信号，服务平滑重启
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		s := <-c
		slog.Info("server stop", "signal", s)
		server.GracefulStop()
	}()

	pb.RegisterDeviceIntServiceServer(server, &deviceapi.DeviceIntService{})
	pb.RegisterMessageExtServiceServer(server, &messageapi.MessageExtService{})
	pb.RegisterMessageIntServiceServer(server, &messageapi.MessageIntService{})
	pb.RegisterGroupIntServiceServer(server, &groupapi.GroupIntService{})
	pb.RegisterRoomIntServiceServer(server, &room.RoomIntService{})

	listen, err := net.Listen("tcp", config.Config.LogicRPCListenAddr)
	if err != nil {
		panic(err)
	}

	slog.Info("rpc服务已经开启")
	err = server.Serve(listen)
	if err != nil {
		slog.Error("serve error", "error", err)
	}
}

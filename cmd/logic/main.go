package main

import (
	"gim/config"
	"gim/internal/logic/api"
	"gim/internal/logic/domain/device"
	"gim/internal/logic/domain/message"
	"gim/internal/logic/proxy"
	"gim/pkg/interceptor"
	"gim/pkg/logger"
	"gim/pkg/protocol/pb"
	"gim/pkg/urlwhitelist"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func init() {
	proxy.MessageProxy = message.App
	proxy.DeviceProxy = device.App
}

func main() {
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.NewInterceptor("logic_interceptor", urlwhitelist.Logic)))

	// 监听服务关闭信号，服务平滑重启
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		s := <-c
		logger.Logger.Info("server stop", zap.Any("signal", s))
		server.GracefulStop()
	}()

	pb.RegisterLogicIntServer(server, &api.LogicIntServer{})
	pb.RegisterLogicExtServer(server, &api.LogicExtServer{})
	listen, err := net.Listen("tcp", config.Config.LogicRPCListenAddr)
	if err != nil {
		panic(err)
	}

	logger.Logger.Info("rpc服务已经开启")
	err = server.Serve(listen)
	if err != nil {
		logger.Logger.Error("serve error", zap.Error(err))
	}
}

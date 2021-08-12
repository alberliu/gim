package main

import (
	"gim/config"
	"gim/internal/logic/api"
	"gim/pkg/db"
	"gim/pkg/interceptor"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc"
	"gim/pkg/urlwhitelist"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger.Init()
	db.InitMysql(config.Logic.MySQL)
	db.InitRedis(config.Logic.RedisIP, config.Logic.RedisPassword)

	// 初始化RpcClient
	rpc.InitConnectIntClient(config.RPCAddr.ConnectRPCAddr)
	rpc.InitBusinessIntClient(config.RPCAddr.BusinessRPCAddr)

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.NewInterceptor("logic_int_interceptor", urlwhitelist.Logic)))

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
	listen, err := net.Listen("tcp", config.Logic.RPCListenAddr)
	if err != nil {
		panic(err)
	}

	logger.Logger.Info("rpc服务已经开启")
	err = server.Serve(listen)
	if err != nil {
		logger.Logger.Error("Serve error", zap.Error(err))
	}
}

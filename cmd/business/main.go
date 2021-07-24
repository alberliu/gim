package main

import (
	"gim/config"
	"gim/internal/business/api"
	"gim/pkg/db"
	"gim/pkg/interceptor"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var whitelistMethed = map[string]int{
	"/pb.BusinessExt/SignIn": 0,
}

func main() {
	logger.Init()
	db.InitMysql(config.Business.MySQL)
	db.InitRedis(config.Business.RedisIP, config.Logic.RedisPassword)

	// 初始化RpcClient
	rpc.InitLogicIntClient(config.Business.LogicRPCAddrs)

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.NewInterceptor("business_interceptor", whitelistMethed)))

	// 监听服务关闭信号，服务平滑重启
	go func() {
		c := make(chan os.Signal, 0)
		signal.Notify(c, syscall.SIGTERM)
		s := <-c
		logger.Logger.Info("server stop", zap.Any("signal", s))
		server.GracefulStop()
	}()

	pb.RegisterBusinessIntServer(server, &api.BusinessIntServer{})
	pb.RegisterBusinessExtServer(server, &api.BusinessExtServer{})
	listen, err := net.Listen("tcp", config.Business.RPCListenAddr)
	if err != nil {
		panic(err)
	}

	logger.Logger.Info("rpc服务已经开启")
	err = server.Serve(listen)
	if err != nil {
		logger.Logger.Error("Serve", zap.Error(err))
	}
}

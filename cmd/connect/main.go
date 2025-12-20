package main

import (
	"google.golang.org/grpc"

	"gim/internal/connect"
	"gim/pkg/logger"
	pb "gim/pkg/protocol/pb/connectpb"
	"gim/pkg/server"
)

func main() {
	logger.Init("connect")

	// 启动TCP长链接服务器
	go func() {
		connect.StartTCPServer(":8002")
	}()

	// 启动WebSocket长链接服务器
	go func() {
		connect.StartWSServer(":8003")
	}()

	// 启动服务订阅
	connect.StartSubscribe()

	server.RunGRPCServer(func(server *grpc.Server) {
		pb.RegisterConnectIntServiceServer(server, &connect.ConnIntService{})
	})

	server.WaitForShutdown()
}

package main

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"

	"gim/internal/connect"
	"gim/pkg/logger"
	pb "gim/pkg/protocol/pb/connectpb"
	"gim/pkg/server"
)

func main() {
	logger.Init()

	// 启动TCP长链接服务器
	tcpListener := connect.StartTCPServer(":8002")

	// 启动WebSocket长链接服务器
	wsServer := connect.StartWSServer(":8003")

	// 启动服务订阅
	connect.StartSubscribe()

	server.RunGRPCServer(func(server *grpc.Server) {
		pb.RegisterConnectIntServiceServer(server, &connect.ConnectIntService{})
	})

	server.WaitForShutdown()

	if err := tcpListener.Close(); err != nil {
		slog.Error("tcp listener close failed", "error", err)
	}
	slog.Info("tcp listener closed")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := wsServer.Shutdown(ctx); err != nil {
		slog.Error("http server shutdown failed", "error", err)
	}
	slog.Info("ws listener closed")
}

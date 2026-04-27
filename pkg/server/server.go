package server

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
	"gim/pkg/ugrpc"
)

var grpcServer *grpc.Server

func RunGRPCServer(f func(server *grpc.Server)) {
	grpcServer = grpc.NewServer(grpc.ChainUnaryInterceptor(ugrpc.ServerInterceptor))
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	f(grpcServer)

	listen, err := net.Listen("tcp", config.GrpcListenAddr)
	if err != nil {
		panic(err)
	}
	go func() {
		slog.Info("rpc server started", "addr", config.GrpcListenAddr)
		err = grpcServer.Serve(listen)
		if err != nil {
			slog.Error("start rpc server failed", "error", err)
			panic(err)
		}
	}()
}

func WaitForShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	s := <-c
	slog.Info("server stopping", "signal", s)

	if grpcServer != nil {
		grpcServer.GracefulStop()
		slog.Info("grpc server stopped")
	}
}

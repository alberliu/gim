package tcp_conn

import (
	"context"
	"gim/config"
	"gim/pkg/gerrors"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"net"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ConnIntServer struct{}

// Message 投递消息
func (s *ConnIntServer) DeliverMessage(ctx context.Context, req *pb.DeliverMessageReq) (*pb.DeliverMessageResp, error) {
	return &pb.DeliverMessageResp{}, DeliverMessage(ctx, req)
}

// UnaryServerInterceptor 服务器端的单向调用的拦截器
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer gerrors.LogPanic("tcp_conn_int_interceptor", ctx, req, info, &err)

	resp, err = handler(ctx, req)
	logger.Logger.Debug("interceptor", zap.Any("info", info), zap.Any("req", req), zap.Any("resp", resp))

	s, _ := status.FromError(err)
	if s.Code() != 0 && s.Code() < 1000 {
		md, _ := metadata.FromIncomingContext(ctx)
		logger.Logger.Error("tcp_conn_int_interceptor", zap.String("method", info.FullMethod), zap.Any("md", md), zap.Any("req", req),
			zap.Any("resp", resp), zap.Error(err), zap.String("stack", gerrors.GetErrorStack(s)))
	}
	return resp, err
}

// StartRPCServer 启动rpc服务器
func StartRPCServer() {
	listener, err := net.Listen("tcp", config.TCPConn.RPCListenAddr)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(UnaryServerInterceptor))
	pb.RegisterConnIntServer(server, &ConnIntServer{})
	logger.Logger.Debug("rpc服务已经开启")
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}

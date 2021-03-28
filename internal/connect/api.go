package connect

import (
	"context"
	"gim/config"
	"gim/pkg/gerrors"
	"gim/pkg/grpclib"
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
	resp := &pb.DeliverMessageResp{}

	// 获取设备对应的TCP连接
	conn := GetConn(req.DeviceId)
	if conn == nil {
		logger.Logger.Warn("GetConn warn", zap.Int64("device_id", req.DeviceId))
		return resp, nil
	}

	if conn.DeviceId != req.DeviceId {
		logger.Logger.Warn("GetConn warn", zap.Int64("device_id", req.DeviceId))
		return resp, nil
	}

	conn.Send(pb.PackageType_PT_MESSAGE, grpclib.GetCtxRequstId(ctx), req.MessageSend, nil)
	return resp, nil
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
	listener, err := net.Listen("tcp", config.Connect.RPCListenAddr)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(UnaryServerInterceptor))
	pb.RegisterConnectIntServer(server, &ConnIntServer{})
	logger.Logger.Debug("rpc服务已经开启")
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}

package tcp_conn

import (
	"context"
	"gim/conf"
	"gim/public/grpclib"
	"gim/public/logger"
	"gim/public/pb"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ConnIntServer struct{}

// Message 投递消息
func (s *ConnIntServer) DeliverMessage(ctx context.Context, req *pb.DeliverMessageReq) (*pb.DeliverMessageResp, error) {
	// 获取设备对应的TCP连接
	conn := load(req.DeviceId)
	if ctx == nil {
		logger.Sugar.Warn("ctx id nil")
		return &pb.DeliverMessageResp{}, nil
	}

	// 发送消息
	conn.Output(pb.PackageType_PT_MESSAGE, grpclib.GetCtxRequstId(ctx), nil, req.Message)
	return &pb.DeliverMessageResp{}, nil
}

// UnaryServerInterceptor 服务器端的单向调用的拦截器
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	logger.Logger.Debug("interceptor", zap.Any("info", info), zap.Any("req", req), zap.Any("resp", resp))
	return resp, err
}

// StartRPCServer 启动rpc服务器
func StartRPCServer() {
	listener, err := net.Listen("tcp", conf.ConnConf.RPCListenAddr)
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

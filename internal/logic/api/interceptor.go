package api

import (
	"context"
	"gim/internal/logic/service"
	"gim/pkg/gerrors"
	"gim/pkg/grpclib"
	"gim/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// LogicIntInterceptor 服务器端的单向调用的拦截器
func LogicIntInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer gerrors.LogPanic("logic_int_interceptor", ctx, req, info, &err)

	resp, err = handler(ctx, req)
	logger.Logger.Debug("logic_int_interceptor", zap.Any("info", info), zap.Any("req", req), zap.Any("resp", resp), zap.Error(err))

	s, _ := status.FromError(err)
	if s.Code() != 0 && s.Code() < 1000 {
		md, _ := metadata.FromIncomingContext(ctx)
		logger.Logger.Error("logic_int_interceptor", zap.String("method", info.FullMethod), zap.Any("md", md), zap.Any("req", req),
			zap.Any("resp", resp), zap.Error(err), zap.String("stack", gerrors.GetErrorStack(s)))
	}
	return resp, err
}

// LogicExtInterceptor 服务器端的单向调用的拦截器
func LogicExtInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer gerrors.LogPanic("logic_ext_interceptor", ctx, req, info, &err)

	resp, err = doLogicClientExt(ctx, req, info, handler)
	logger.Logger.Debug("logic_ext_interceptor", zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),
		zap.Any("resp", resp), zap.Error(err))

	s, _ := status.FromError(err)
	if s.Code() != 0 && s.Code() < 1000 {
		md, _ := metadata.FromIncomingContext(ctx)
		logger.Logger.Error("logic_ext_interceptor", zap.String("method", info.FullMethod), zap.Any("md", md), zap.Any("req", req),
			zap.Any("resp", resp), zap.Error(err), zap.String("stack", gerrors.GetErrorStack(s)))
	}
	return
}

func doLogicClientExt(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod != "/pb.LogicExt/RegisterDevice" {
		userId, deviceId, err := grpclib.GetCtxData(ctx)
		if err != nil {
			return nil, err
		}
		token, err := grpclib.GetCtxToken(ctx)
		if err != nil {
			return nil, err
		}

		err = service.AuthService.Auth(ctx, userId, deviceId, token)
		if err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

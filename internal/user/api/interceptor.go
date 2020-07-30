package api

import (
	"context"
	"gim/internal/user/service"
	"gim/pkg/gerrors"
	"gim/pkg/grpclib"
	"gim/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	IntServerName = "user_int"
	ExtServerName = "user_ext"
)

// 服务器端的单向调用的拦截器
func UserIntInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer gerrors.LogPanic(IntServerName, ctx, req, info, &err)

	resp, err = handler(ctx, req)
	logger.Logger.Debug(IntServerName, zap.Any("info", info), zap.Any("req", req), zap.Any("resp", resp), zap.Error(err))

	s, _ := status.FromError(err)
	if s.Code() != 0 && s.Code() < 1000 {
		md, _ := metadata.FromIncomingContext(ctx)
		logger.Logger.Error(IntServerName, zap.String("method", info.FullMethod), zap.Any("md", md), zap.Any("req", req),
			zap.Any("resp", resp), zap.Error(err), zap.String("stack", gerrors.GetErrorStack(s)))
	}
	return resp, err
}

// 服务器端的单向调用的拦截器
func UserExtInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer gerrors.LogPanic(ExtServerName, ctx, req, info, &err)

	resp, err = doLogicExt(ctx, req, info, handler)
	logger.Logger.Debug(ExtServerName, zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),
		zap.Any("resp", resp), zap.Error(err))

	s, _ := status.FromError(err)
	if s.Code() != 0 && s.Code() < 1000 {
		md, _ := metadata.FromIncomingContext(ctx)
		logger.Logger.Error(ExtServerName, zap.String("method", info.FullMethod), zap.Any("md", md), zap.Any("req", req),
			zap.Any("resp", resp), zap.Error(err), zap.String("stack", gerrors.GetErrorStack(s)))
	}
	return
}

func doLogicExt(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod != "/pb.UserExt/SignIn" {
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

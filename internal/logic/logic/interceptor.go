package logic

import (
	"context"
	"gim/internal/logic/service"
	"gim/pkg/gerrors"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/util"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func logPanic(serverName string, ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, err *error) {
	p := recover()
	if p != nil {
		logger.Logger.Error(serverName+" panic", zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),
			zap.Any("panic", p), zap.String("stack", util.GetStackInfo()))
		*err = gerrors.ErrUnknown
	}
}

// 服务器端的单向调用的拦截器
func LogicIntInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		logPanic("logic_int_interceptor", ctx, req, info, &err)
	}()

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

// 服务器端的单向调用的拦截器
func LogicClientExtInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		logPanic("logic_client_ext_interceptor", ctx, req, info, &err)
	}()

	resp, err = doLogicClientExt(ctx, req, info, handler)
	logger.Logger.Debug("logic_client_ext_interceptor", zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),
		zap.Any("resp", resp), zap.Error(err))

	s, _ := status.FromError(err)
	if s.Code() != 0 && s.Code() < 1000 {
		md, _ := metadata.FromIncomingContext(ctx)
		logger.Logger.Error("logic_client_ext_interceptor", zap.String("method", info.FullMethod), zap.Any("md", md), zap.Any("req", req),
			zap.Any("resp", resp), zap.Error(err), zap.String("stack", gerrors.GetErrorStack(s)))
	}
	return
}

func doLogicClientExt(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod != "/pb.LogicClientExt/RegisterDevice" {
		appId, userId, deviceId, err := grpclib.GetCtxData(ctx)
		if err != nil {
			return nil, err
		}
		token, err := grpclib.GetCtxToken(ctx)
		if err != nil {
			return nil, err
		}

		err = service.AuthService.Auth(ctx, appId, userId, deviceId, token)
		if err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

// 服务器端的单向调用的拦截器
func LogicServerExtInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		logPanic("logic_server_ext_interceptor", ctx, req, info, &err)
	}()

	resp, err = doLogicServerExt(ctx, req, info, handler)
	logger.Logger.Debug("logic_server_ext_interceptor", zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),
		zap.Any("resp", resp), zap.Error(err))

	s, _ := status.FromError(err)
	if s.Code() != 0 && s.Code() < 1000 {
		md, _ := metadata.FromIncomingContext(ctx)
		logger.Logger.Error("logic_server_ext_interceptor", zap.String("method", info.FullMethod), zap.Any("md", md), zap.Any("req", req),
			zap.Any("resp", resp), zap.Error(err), zap.String("stack", gerrors.GetErrorStack(s)))
	}
	return resp, err
}

func doLogicServerExt(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	appId, err := grpclib.GetCtxAppId(ctx)
	if err != nil {
		return nil, err
	}
	token, err := grpclib.GetCtxToken(ctx)
	if err != nil {
		return nil, err
	}

	err = service.AuthService.Auth(ctx, appId, 0, 0, token)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

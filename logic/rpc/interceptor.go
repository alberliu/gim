package rpc

import (
	"context"
	"gim/logic/service"
	"gim/public/grpclib"
	"gim/public/imerror"
	"gim/public/logger"
	"gim/public/util"

	"google.golang.org/grpc/status"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// 服务器端的单向调用的拦截器
func LogicIntInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	logger.Logger.Debug("logic_int_interceptor", zap.Any("info", info), zap.Any("req", req), zap.Any("resp", resp), zap.Error(err))

	if _, ok := status.FromError(err); !ok {
		logger.Logger.Error("logic_int_interceptor", zap.Any("info", info), zap.Any("req", req), zap.Any("resp", resp), zap.Error(err))
	}
	return resp, err
}

// 服务器端的单向调用的拦截器
func LogicClientExtInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		p := recover()
		if p != nil {
			logger.Logger.Debug("logic_client_ext_interceptor panic", zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),
				zap.Any("panic", p), zap.String("stack", util.GetStackInfo()))
			err = imerror.ErrUnknown
		}
	}()

	resp, err = doLogicClientExt(ctx, req, info, handler)
	logger.Logger.Debug("logic_client_ext_interceptor", zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),
		zap.Any("resp", resp), zap.Error(err))
	if _, ok := status.FromError(err); !ok {
		logger.Logger.Error("logic_client_ext_interceptor", zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),
			zap.Any("resp", resp), zap.Error(err))
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

		err = service.AuthService.Auth(Context(), appId, userId, deviceId, token)
		if err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

// 服务器端的单向调用的拦截器
func LogicServerExtInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		p := recover()
		if p != nil {
			logger.Logger.Debug("logic_server_ext_interceptor panic", zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),
				zap.Any("panic", p), zap.String("stack", util.GetStackInfo()))
			err = imerror.ErrUnknown
		}
	}()

	resp, err = doLogicServerExt(ctx, req, info, handler)
	logger.Logger.Debug("logic_server_ext_interceptor", zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),
		zap.Any("resp", resp), zap.Error(err))
	if _, ok := status.FromError(err); !ok {
		logger.Logger.Error("logic_server_ext_interceptor", zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),
			zap.Any("resp", resp), zap.Error(err))
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

	err = service.AuthService.Auth(Context(), appId, 0, 0, token)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

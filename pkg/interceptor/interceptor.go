package interceptor

import (
	"context"
	"log/slog"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"gim/pkg/gerrors"
	"gim/pkg/md"
	"gim/pkg/protocol/pb/userpb"
	"gim/pkg/rpc"
)

// NewInterceptor 生成GRPC过滤器
func NewInterceptor(urlWhitelist map[string]struct{}) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
		defer gerrors.LogPanic(ctx, req, info, &err)
		md, _ := metadata.FromIncomingContext(ctx)
		logger := slog.With("method", info.FullMethod, "md", md, "request", req, "reply", reply)

		reply, err = handleWithAuth(ctx, req, info, handler, urlWhitelist)

		s, _ := status.FromError(err)
		if s.Code() != 0 && s.Code() < 10000 {
			logger.Error("interceptor", "error", err, "stack", gerrors.GetErrorStack(s))
		}
		logger.Debug("interceptor", "error", err)
		return
	}
}

// handleWithAuth 处理鉴权逻辑
func handleWithAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, urlWhitelist map[string]struct{}) (interface{}, error) {
	serverName := strings.Split(info.FullMethod, "/")[1]
	if !strings.HasSuffix(serverName, "IntService") {
		if _, ok := urlWhitelist[info.FullMethod]; !ok {
			userID, deviceID, err := md.GetCtxData(ctx)
			if err != nil {
				return nil, err
			}
			token := md.GetCtxToken(ctx)

			_, err = rpc.GetUserIntClient().Auth(ctx, &userpb.AuthRequest{
				UserId:   userID,
				DeviceId: deviceID,
				Token:    token,
			})

			if err != nil {
				return nil, err
			}
		}
	}
	return handler(ctx, req)
}

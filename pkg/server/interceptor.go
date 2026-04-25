package server

import (
	"context"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"gim/pkg/gerrors"
	"gim/pkg/logger"
	"gim/pkg/md"
	"gim/pkg/protocol/pb/businesspb"
	"gim/pkg/rpc"
)

// NewInterceptor 生成GRPC过滤器
func NewInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply any, err error) {
		defer gerrors.LogPanic(ctx, req, info, &err)

		ctx = generateRequestID(ctx)
		reply, err = handleWithAuth(ctx, req, info, handler)

		inMD, _ := metadata.FromIncomingContext(ctx)
		log := slog.With("method", info.FullMethod, "md", inMD, "request", req, "reply", reply)

		s, _ := status.FromError(err)
		if s.Code() != 0 && s.Code() < gerrors.MinBizCode {
			log.ErrorContext(ctx, "handle request failed", logger.Error(err), "stack", gerrors.GetErrorStack(s))
		}
		log.DebugContext(ctx, "handle request", logger.Error(err))
		return
	}
}

func generateRequestID(ctx context.Context) context.Context {
	requestID := md.GetRequestID(ctx)
	if requestID == "" {
		requestID = uuid.NewString()
		incomingMD, _ := metadata.FromIncomingContext(ctx)
		incomingMD = incomingMD.Copy()
		incomingMD.Set(md.RequestID, requestID)
		ctx = metadata.NewIncomingContext(ctx, incomingMD)
	}
	_ = grpc.SetHeader(ctx, metadata.Pairs(md.RequestID, requestID))
	return ctx
}

// handleWithAuth 处理鉴权逻辑
func handleWithAuth(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	serverName := strings.Split(info.FullMethod, "/")[1]
	if strings.HasSuffix(serverName, "ExtService") {
		if _, ok := URLWhitelist[info.FullMethod]; !ok {
			userID := md.GetUserID(ctx)
			token := md.GetToken(ctx)

			_, err := rpc.GetUserIntClient().Auth(ctx, &businesspb.AuthRequest{
				UserId: userID,
				Token:  token,
			})

			if err != nil {
				return nil, err
			}
		}
	}
	return handler(ctx, req)
}

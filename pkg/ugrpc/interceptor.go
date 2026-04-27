package ugrpc

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"gim/pkg/gerrors"
	"gim/pkg/logger"
	"gim/pkg/md"
)

func ProxyLogClientInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx = generateRequestID(ctx)

	inMD, _ := metadata.FromIncomingContext(ctx)
	ctx = metadata.NewOutgoingContext(ctx, inMD)
	err := invoker(ctx, method, req, reply, cc, opts...)

	log := slog.With("method", method, "md", inMD)
	s, _ := status.FromError(err)
	if s.Code() != 0 && s.Code() < gerrors.MinBizCode {
		log.ErrorContext(ctx, "send request failed", logger.Error(err), "stack", gerrors.GetErrorStack(s))
	}
	log.DebugContext(ctx, "send request", logger.Error(err))
	return err
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

func ServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply any, err error) {
	defer gerrors.LogPanic(ctx, req, info, &err)

	inMD, _ := metadata.FromIncomingContext(ctx)
	ctx = metadata.NewOutgoingContext(ctx, inMD)
	reply, err = handler(ctx, req)

	log := slog.With("method", info.FullMethod, "md", inMD, "request", req, "reply", reply)
	s, _ := status.FromError(err)
	if s.Code() != 0 && s.Code() < gerrors.MinBizCode {
		log.ErrorContext(ctx, "handle request failed", logger.Error(err), "stack", gerrors.GetErrorStack(s))
	}
	log.DebugContext(ctx, "handle request", logger.Error(err))
	return
}

func clientInterceptor(ctx context.Context, method string, request, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	ctx = md.WithRequestID(ctx, md.GetRequestID(ctx))
	err := invoker(ctx, method, request, reply, cc, opts...)

	md, _ := metadata.FromOutgoingContext(ctx)
	slog.DebugContext(ctx, "send request", "method", method, "metadata", md, "request", request, "reply", reply, logger.Error(err))
	return err
}

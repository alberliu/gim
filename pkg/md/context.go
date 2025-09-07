package md

import (
	"context"
	"strconv"

	"google.golang.org/grpc/metadata"
)

const (
	CtxUserID    = "user_id"
	CtxDeviceID  = "device_id"
	CtxToken     = "token"
	CtxRequestID = "request_id"
)

func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(CtxRequestID, requestID))
}

func Get(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values, ok := md[key]
	if !ok || len(values) == 0 {
		return ""
	}
	return values[0]
}

// GetRequestID 获取ctx的app_id
func GetRequestID(ctx context.Context) string {
	return Get(ctx, CtxRequestID)
}

func GetUserID(ctx context.Context) uint64 {
	userID, _ := strconv.ParseUint(Get(ctx, CtxUserID), 10, 64)
	return userID
}

func GetDeviceID(ctx context.Context) uint64 {
	deviceID, _ := strconv.ParseUint(Get(ctx, CtxDeviceID), 10, 64)
	return deviceID
}

// GetToken 获取ctx的token
func GetToken(ctx context.Context) string {
	return Get(ctx, CtxToken)
}

// NewAndCopyRequestID 创建一个context,并且复制RequestID
func NewAndCopyRequestID(ctx context.Context) context.Context {
	newCtx := context.TODO()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return newCtx
	}

	requestIDs, ok := md[CtxRequestID]
	if !ok && len(requestIDs) == 0 {
		return newCtx
	}
	return metadata.NewOutgoingContext(newCtx, metadata.Pairs(CtxRequestID, requestIDs[0]))
}

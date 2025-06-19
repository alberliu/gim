package md

import (
	"context"
	"strconv"

	"google.golang.org/grpc/metadata"

	"gim/pkg/gerrors"
)

const (
	CtxUserID    = "user_id"
	CtxDeviceID  = "device_id"
	CtxToken     = "token"
	CtxRequestID = "request_id"
)

func ContextWithRequestID(ctx context.Context, requestID int64) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(CtxRequestID, strconv.FormatInt(requestID, 10)))
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

// GetCtxRequestID 获取ctx的app_id
func GetCtxRequestID(ctx context.Context) int64 {
	requestIDStr := Get(ctx, CtxRequestID)
	requestID, err := strconv.ParseInt(requestIDStr, 10, 64)
	if err != nil {
		return 0
	}
	return requestID
}

// GetCtxData 获取ctx的用户数据，依次返回user_id,device_id
func GetCtxData(ctx context.Context) (uint64, uint64, error) {
	var (
		userID   uint64
		deviceID uint64
		err      error
	)

	userIDStr := Get(ctx, CtxUserID)
	userID, err = strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return 0, 0, gerrors.ErrUnauthorized
	}

	deviceIDStr := Get(ctx, CtxDeviceID)
	deviceID, err = strconv.ParseUint(deviceIDStr, 10, 64)
	if err != nil {
		return 0, 0, gerrors.ErrUnauthorized
	}
	return userID, deviceID, nil
}

// GetCtxToken 获取ctx的token
func GetCtxToken(ctx context.Context) string {
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

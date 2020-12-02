package grpclib

import (
	"context"
	"gim/pkg/gerrors"
	"gim/pkg/logger"
	"strconv"

	"google.golang.org/grpc/metadata"
)

const (
	CtxUserId    = "user_id"
	CtxDeviceId  = "device_id"
	CtxToken     = "token"
	CtxRequestId = "request_id"
)

func ContextWithRequstId(ctx context.Context, requestId int64) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(CtxRequestId, strconv.FormatInt(requestId, 10)))
}

// GetCtxAppId 获取ctx的app_id
func GetCtxRequstId(ctx context.Context) int64 {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0
	}

	requstIds, ok := md[CtxRequestId]
	if !ok && len(requstIds) == 0 {
		return 0
	}
	requstId, err := strconv.ParseInt(requstIds[0], 10, 64)
	if err != nil {
		return 0
	}
	return requstId
}

// GetCtxData 获取ctx的用户数据，依次返回user_id,device_id
func GetCtxData(ctx context.Context) (int64, int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, 0, gerrors.ErrUnauthorized
	}

	var (
		userId   int64
		deviceId int64
		err      error
	)

	userIdStrs, ok := md[CtxUserId]
	if !ok && len(userIdStrs) == 0 {
		return 0, 0, gerrors.ErrUnauthorized
	}
	userId, err = strconv.ParseInt(userIdStrs[0], 10, 64)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, 0, gerrors.ErrUnauthorized
	}

	deviceIdStrs, ok := md[CtxDeviceId]
	if !ok && len(deviceIdStrs) == 0 {
		return 0, 0, gerrors.ErrUnauthorized
	}
	deviceId, err = strconv.ParseInt(deviceIdStrs[0], 10, 64)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, 0, gerrors.ErrUnauthorized
	}
	return userId, deviceId, nil
}

// GetCtxDeviceId 获取ctx的设备id
func GetCtxDeviceId(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, gerrors.ErrUnauthorized
	}

	deviceIdStrs, ok := md[CtxDeviceId]
	if !ok && len(deviceIdStrs) == 0 {
		return 0, gerrors.ErrUnauthorized
	}
	deviceId, err := strconv.ParseInt(deviceIdStrs[0], 10, 64)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, gerrors.ErrUnauthorized
	}
	return deviceId, nil
}

// GetCtxToken 获取ctx的token
func GetCtxToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", gerrors.ErrUnauthorized
	}

	tokens, ok := md[CtxToken]
	if !ok && len(tokens) == 0 {
		return "", gerrors.ErrUnauthorized
	}

	return tokens[0], nil
}

// NewAndCopyRequestId 创建一个context,并且复制RequestId
func NewAndCopyRequestId(ctx context.Context) context.Context {
	newCtx := context.TODO()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return newCtx
	}

	requstIds, ok := md[CtxRequestId]
	if !ok && len(requstIds) == 0 {
		return newCtx
	}
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(CtxRequestId, requstIds[0]))
}

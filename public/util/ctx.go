package util

import (
	"context"
	"gim/public/imerror"
	"gim/public/logger"
	"google.golang.org/grpc/metadata"
	"strconv"
)

const (
	CtxAppId    = "app_id"
	CtxUserId   = "user_id"
	CtxDeviceId = "device_id"
)

func GetCtxData(ctx context.Context) (int64, int64, int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, 0, 0, imerror.ErrUnauthorized
	}

	var (
		appId    int64
		userId   int64
		deviceId int64
		err      error
	)

	// app_id是必填项
	appIdStrs, ok := md[CtxAppId]
	if !ok && len(appIdStrs) == 0 {
		return 0, 0, 0, imerror.ErrUnauthorized
	}
	appId, err = strconv.ParseInt(appIdStrs[0], 10, 64)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, 0, 0, imerror.ErrUnauthorized
	}

	userIdStrs, ok := md[CtxUserId]
	if !ok && len(userIdStrs) == 0 {
		return 0, 0, 0, imerror.ErrUnauthorized
	}
	userId, err = strconv.ParseInt(userIdStrs[0], 10, 64)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, 0, 0, imerror.ErrUnauthorized
	}

	deviceIdStrs, ok := md[CtxDeviceId]
	if !ok && len(deviceIdStrs) == 0 {
		return 0, 0, 0, imerror.ErrUnauthorized
	}
	deviceId, err = strconv.ParseInt(deviceIdStrs[0], 10, 64)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, 0, 0, imerror.ErrUnauthorized
	}
	return appId, userId, deviceId, nil
}

func GetCtxToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", imerror.ErrUnauthorized
	}

	tokens, ok := md[CtxAppId]
	if !ok && len(tokens) == 0 {
		return "", imerror.ErrUnauthorized
	}

	return tokens[0], nil
}

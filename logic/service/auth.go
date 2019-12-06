package service

import (
	"gim/public/imctx"
	"gim/public/imerror"
	"gim/public/logger"
	"gim/public/util"
	"time"
)

type authService struct{}

var AuthService = new(authService)

// SignIn 长连接登录
func (*authService) SignIn(ctx *imctx.Context, appId, userId, deviceId int64, token string, connectAddr string) error {
	// 用户验证
	err := AuthService.VerifyToken(ctx, appId, userId, deviceId, token)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	// 标记用户在设备上登录
	err = DeviceService.Online(ctx, appId, deviceId, userId, connectAddr)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	return nil
}

// Auth 验证用户是否登录
func (*authService) Auth(ctx *imctx.Context, appId, userId, deviceId int64, token string) error {
	return AuthService.VerifyToken(ctx, appId, userId, deviceId, token)
}

// VerifySecretKey 对用户秘钥进行校验
func (*authService) VerifyToken(ctx *imctx.Context, appId, userId, deviceId int64, token string) error {
	app, err := AppService.Get(ctx, appId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	if app == nil {
		return imerror.ErrBadRequest
	}

	info, err := util.DecryptToken(token, app.PrivateKey)
	if err != nil {
		logger.Sugar.Error(err)
		return imerror.ErrUnauthorized
	}

	if !(info.AppId == appId && info.UserId == userId && info.DeviceId == deviceId) {
		return imerror.ErrUnauthorized
	}

	if info.Expire < time.Now().Unix() {
		return imerror.ErrUnauthorized
	}
	return nil
}

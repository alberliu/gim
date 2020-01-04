package service

import (
	"context"
	"gim/pkg/gerrors"
	"gim/pkg/util"
	"time"
)

type authService struct{}

var AuthService = new(authService)

// SignIn 长连接登录
func (*authService) SignIn(ctx context.Context, appId, userId, deviceId int64, token string, connectAddr string) error {
	// 用户验证
	err := AuthService.VerifyToken(ctx, appId, userId, deviceId, token)
	if err != nil {
		return err
	}

	// 标记用户在设备上登录
	err = DeviceService.Online(ctx, appId, deviceId, userId, connectAddr)
	if err != nil {
		return err
	}

	return nil
}

// Auth 验证用户是否登录
func (*authService) Auth(ctx context.Context, appId, userId, deviceId int64, token string) error {
	return AuthService.VerifyToken(ctx, appId, userId, deviceId, token)
}

// VerifySecretKey 对用户秘钥进行校验
func (*authService) VerifyToken(ctx context.Context, appId, userId, deviceId int64, token string) error {
	app, err := AppService.Get(ctx, appId)
	if err != nil {
		return err
	}

	if app == nil {
		return gerrors.ErrBadRequest
	}

	info, err := util.DecryptToken(token, app.PrivateKey)
	if err != nil {
		return gerrors.ErrUnauthorized
	}

	if !(info.AppId == appId && info.UserId == userId && info.DeviceId == deviceId) {
		return gerrors.ErrUnauthorized
	}

	if info.Expire < time.Now().Unix() {
		return gerrors.ErrUnauthorized
	}
	return nil
}

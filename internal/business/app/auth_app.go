package app

import (
	"context"
	"gim/internal/business/domain/user/service"
)

type authApp struct{}

var AuthApp = new(authApp)

// SignIn 长连接登录
func (*authApp) SignIn(ctx context.Context, phoneNumber, code string, deviceId int64) (bool, int64, string, error) {
	return service.AuthService.SignIn(ctx, phoneNumber, code, deviceId)
}

// Auth 验证用户是否登录
func (*authApp) Auth(ctx context.Context, userId, deviceId int64, token string) error {
	return service.AuthService.Auth(ctx, userId, deviceId, token)
}

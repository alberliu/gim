package service

import (
	"context"
	"gim/pkg/pb"
	"gim/pkg/rpc"
)

type authService struct{}

var AuthService = new(authService)

// SignIn 长连接登录
func (*authService) SignIn(ctx context.Context, userId, deviceId int64, token string, connAddr string, connFd int64) error {
	_, err := rpc.UserIntClient.Auth(ctx, &pb.AuthReq{UserId: userId, DeviceId: deviceId, Token: token})
	if err != nil {
		return err
	}

	// 标记用户在设备上登录
	err = DeviceService.Online(ctx, deviceId, userId, connAddr, connFd)
	if err != nil {
		return err
	}
	return nil
}

// Auth 权限验证
func (*authService) Auth(ctx context.Context, userId, deviceId int64, token string) error {
	_, err := rpc.UserIntClient.Auth(ctx, &pb.AuthReq{UserId: userId, DeviceId: deviceId, Token: token})
	if err != nil {
		return err
	}
	return nil
}

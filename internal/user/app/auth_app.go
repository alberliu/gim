package app

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis"

	"gim/internal/user/domain"
	"gim/internal/user/repo"
	"gim/pkg/gerrors"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/rpc"
)

type authApp struct{}

var AuthApp = new(authApp)

// SignIn 长连接登录
func (*authApp) SignIn(ctx context.Context, phoneNumber, code string, deviceId uint64) (bool, uint64, string, error) {
	if !verify(phoneNumber, code) {
		return false, 0, "", gerrors.ErrBadCode
	}

	user, err := repo.UserRepo.GetByPhoneNumber(phoneNumber)
	if err != nil && !errors.Is(err, gerrors.ErrUserNotFound) {
		return false, 0, "", err
	}

	var isNew = false
	if errors.Is(err, gerrors.ErrUserNotFound) {
		user = &domain.User{
			PhoneNumber: phoneNumber,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := repo.UserRepo.Save(user)
		if err != nil {
			return false, 0, "", err
		}
		isNew = true
	}

	resp, err := rpc.GetDeviceIntClient().GetDevice(ctx, &logicpb.GetDeviceRequest{DeviceId: deviceId})
	if err != nil {
		return false, 0, "", err
	}

	// 方便测试
	token := "0"
	//token := util.RandString(40)
	err = repo.AuthRepo.Set(user.Id, resp.Device.DeviceId, domain.Device{
		Type:   resp.Device.Type,
		Token:  token,
		Expire: time.Now().AddDate(0, 3, 0).Unix(),
	})
	if err != nil {
		return false, 0, "", err
	}

	return isNew, user.Id, token, nil
}

func verify(phoneNumber, code string) bool {
	// 假装他成功了
	return true
}

// Auth 验证用户是否登录
func (*authApp) Auth(ctx context.Context, userId, deviceId uint64, token string) error {
	device, err := repo.AuthRepo.Get(userId, deviceId)
	if errors.Is(err, redis.Nil) {
		return gerrors.ErrUnauthorized
	}
	if err != nil {
		return err
	}

	if device.Expire < time.Now().Unix() {
		return gerrors.ErrUnauthorized
	}
	if device.Token != token {
		return gerrors.ErrUnauthorized
	}
	return nil
}

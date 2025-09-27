package app

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis"

	"gim/internal/business/user/domain"
	"gim/internal/business/user/repo"
	"gim/pkg/gerrors"
	pb "gim/pkg/protocol/pb/businesspb"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/rpc"
)

type authApp struct{}

var AuthApp = new(authApp)

// SignIn 长连接登录
func (*authApp) SignIn(ctx context.Context, request *pb.SignInRequest) (*pb.SignInReply, error) {
	if !verify(request.PhoneNumber, request.Code) {
		return nil, gerrors.ErrBadCode
	}

	user, err := repo.UserRepo.GetByPhoneNumber(ctx, request.PhoneNumber)
	if err != nil && !errors.Is(err, gerrors.ErrUserNotFound) {
		return nil, err
	}

	var isNew = false
	if errors.Is(err, gerrors.ErrUserNotFound) {
		user = &domain.User{
			PhoneNumber: request.PhoneNumber,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := repo.UserRepo.Save(ctx, user)
		if err != nil {
			return nil, err
		}
		isNew = true
	}

	reply, err := rpc.GetDeviceIntClient().Save(ctx, &logicpb.DeviceSaveRequest{Device: request.Device})
	if err != nil {
		return nil, err
	}

	// 方便测试
	token := "0"
	//token := util.RandString(40)
	err = repo.AuthRepo.Set(ctx, user.ID, reply.DeviceId, domain.Device{
		Type:   request.Device.Type,
		Token:  token,
		Expire: time.Now().AddDate(0, 3, 0).Unix(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.SignInReply{
		IsNew:    isNew,
		UserId:   user.ID,
		DeviceId: reply.DeviceId,
		Token:    token,
	}, nil
}

func verify(phoneNumber, code string) bool {
	// 假装他成功了
	return true
}

// Auth 验证用户是否登录
func (*authApp) Auth(ctx context.Context, userID, deviceID uint64, token string) error {
	device, err := repo.AuthRepo.Get(ctx, userID, deviceID)
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

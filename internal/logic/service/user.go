package service

import (
	"context"
	"gim/internal/logic/cache"
	"gim/internal/logic/dao"
	"gim/internal/logic/model"
	"gim/pkg/gerrors"
)

type userService struct{}

var UserService = new(userService)

// Add 添加用户（将业务账号导入IM系统账户）
//1.添加用户，2.添加用户消息序列号
func (*userService) Add(ctx context.Context, user model.User) error {
	affected, err := dao.UserDao.Add(user)
	if err != nil {
		return err
	}
	if affected == 0 {
		return gerrors.ErrUserAlreadyExist
	}
	return nil
}

// Get 获取用户信息
func (*userService) Get(ctx context.Context, appId, userId int64) (*model.User, error) {
	user, err := cache.UserCache.Get(appId, userId)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	user, err = dao.UserDao.Get(appId, userId)
	if err != nil {
		return nil, err
	}

	if user != nil {
		err = cache.UserCache.Set(*user)
		if err != nil {
			return nil, err
		}
	}
	return user, err
}

// Get 获取用户信息
func (*userService) Update(ctx context.Context, user model.User) error {
	err := dao.UserDao.Update(user)
	if err != nil {
		return err
	}

	err = cache.UserCache.Del(user.AppId, user.UserId)
	if err != nil {
		return err
	}

	return nil
}

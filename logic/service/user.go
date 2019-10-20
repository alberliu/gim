package service

import (
	"goim/logic/cache"
	"goim/logic/dao"
	"goim/logic/model"
	"goim/public/imctx"
	"goim/public/imerror"
	"goim/public/logger"
)

type userService struct{}

var UserService = new(userService)

// Add 添加用户（将业务账号导入IM系统账户）
//1.添加用户，2.添加用户消息序列号
func (*userService) Add(ctx *imctx.Context, user model.User) error {
	affected, err := dao.UserDao.Add(ctx, user)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	if affected == 0 {
		return imerror.ErrUserAlreadyExist
	}
	return nil
}

// Get 获取用户信息
func (*userService) Get(ctx *imctx.Context, appId, userId int64) (*model.User, error) {
	user, err := cache.UserCache.Get(appId, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	user, err = dao.UserDao.Get(ctx, appId, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	if user != nil {
		err = cache.UserCache.Set(*user)
		if err != nil {
			logger.Sugar.Error(err)
			return nil, err
		}
	}
	return user, err
}

// Get 获取用户信息
func (*userService) Update(ctx *imctx.Context, user model.User) error {
	err := dao.UserDao.Update(ctx, user)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	err = cache.UserCache.Del(user.AppId, user.UserId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	return nil
}

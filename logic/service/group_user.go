package service

import (
	"goim/logic/cache"
	"goim/logic/dao"
	"goim/logic/model"
	"goim/public/imctx"
	"goim/public/logger"
)

type groupUserService struct{}

var GroupUserService = new(groupUserService)

// ListByUserId 获取用户所加入的群组
func (*groupUserService) ListByUserId(ctx *imctx.Context, appId, userId int64) ([]model.Group, error) {
	groups, err := dao.GroupUserDao.ListByUserId(ctx, appId, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return groups, nil
}

// GetUsers 获取群组的所有用户信息
func (*groupUserService) GetUsers(ctx *imctx.Context, appId, groupId int64) ([]model.GroupUser, error) {
	users, err := cache.GroupUserCache.Get(appId, groupId)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	if users != nil {
		return users, nil
	}

	users, err = dao.GroupUserDao.ListUser(ctx, appId, groupId)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	err = cache.GroupUserCache.Set(appId, groupId, users)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return users, err
}

// AddUser 给群组添加用户
func (*groupUserService) AddUser(ctx *imctx.Context, appId, groupId, userId int64, label, extra string) error {
	err := dao.GroupUserDao.Add(ctx, appId, groupId, userId, label, extra)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	err = dao.GroupDao.UpdateUserNum(ctx, appId, groupId, 1)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	err = cache.GroupUserCache.Del(appId, groupId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	return nil
}

// DeleteUser 从群组移除用户
func (*groupUserService) DeleteUser(ctx *imctx.Context, appId, groupId, userId int64) error {
	err := dao.GroupUserDao.Delete(ctx, appId, groupId, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	err = dao.GroupDao.UpdateUserNum(ctx, appId, groupId, -1)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	err = cache.GroupUserCache.Del(appId, groupId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	return nil
}

// Update 更新群组用户信息
func (*groupUserService) Update(ctx *imctx.Context, appId, groupId int64, userId int64, label, extra string) error {
	err := dao.GroupUserDao.Update(ctx, appId, groupId, userId, label, extra)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	err = cache.GroupUserCache.Del(appId, groupId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

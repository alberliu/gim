package service

import (
	"gim/logic/cache"
	"gim/logic/dao"
	"gim/logic/model"
	"gim/public/imctx"
	"gim/public/imerror"
	"gim/public/logger"
)

type groupService struct{}

var GroupService = new(groupService)

// Get 获取群组信息
func (*groupService) Get(ctx *imctx.Context, appId, groupId int64) (*model.Group, error) {
	group, err := cache.GroupCache.Get(appId, groupId)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	if group != nil {
		return group, nil
	}
	group, err = dao.GroupDao.Get(ctx, appId, groupId)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return group, nil
}

// Create 创建群组
func (*groupService) Create(ctx *imctx.Context, group model.Group) error {
	affected, err := dao.GroupDao.Add(ctx, group)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	if affected == 0 {
		return imerror.ErrGroupAlreadyExist
	}
	return nil
}

// Update 更新群组
func (*groupService) Update(ctx *imctx.Context, group model.Group) error {
	err := dao.GroupDao.Update(ctx, group.AppId, group.GroupId, group.Name, group.Introduction, group.Extra)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	err = cache.GroupCache.Del(group.AppId, group.GroupId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// AddUser 给群组添加用户
func (*groupService) AddUser(ctx *imctx.Context, appId, groupId, userId int64, label, extra string) error {
	group, err := GroupService.Get(ctx, appId, groupId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	if group == nil {
		return imerror.ErrGroupNotExist
	}

	user, err := UserService.Get(ctx, appId, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	if user == nil {
		return imerror.ErrUserNotExist
	}

	if group.Type == model.GroupTypeGroup {
		err = GroupUserService.AddUser(ctx, appId, groupId, userId, label, extra)
		if err != nil {
			logger.Sugar.Error(err)
			return err
		}
	}
	if group.Type == model.GroupTypeChatRoom {
		err = cache.LargeGroupUserCache.Set(appId, groupId, userId, label, extra)
		if err != nil {
			logger.Sugar.Error(err)
			return err
		}
	}
	return nil
}

// UpdateUser 更新群组用户
func (*groupService) UpdateUser(ctx *imctx.Context, appId, groupId, userId int64, label, extra string) error {
	group, err := GroupService.Get(ctx, appId, groupId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	if group == nil {
		return imerror.ErrGroupNotExist
	}

	if group.Type == model.GroupTypeGroup {
		err = GroupUserService.Update(ctx, appId, groupId, userId, label, extra)
		if err != nil {
			logger.Sugar.Error(err)
			return err
		}
	}
	if group.Type == model.GroupTypeChatRoom {
		err = cache.LargeGroupUserCache.Set(appId, groupId, userId, label, extra)
		if err != nil {
			logger.Sugar.Error(err)
			return err
		}
	}
	return nil
}

// DeleteUser 删除用户群组
func (*groupService) DeleteUser(ctx *imctx.Context, appId, groupId, userId int64) error {
	group, err := GroupService.Get(ctx, appId, groupId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	if group == nil {
		return imerror.ErrGroupNotExist
	}

	if group.Type == model.GroupTypeGroup {
		err = GroupUserService.DeleteUser(ctx, appId, groupId, userId)
		if err != nil {
			logger.Sugar.Error(err)
			return err
		}
	}
	if group.Type == model.GroupTypeChatRoom {
		err = cache.LargeGroupUserCache.Del(appId, groupId, userId)
		if err != nil {
			logger.Sugar.Error(err)
			return err
		}
	}
	return nil
}

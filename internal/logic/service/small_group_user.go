package service

import (
	"context"
	"gim/internal/logic/cache"
	"gim/internal/logic/dao"
	"gim/internal/logic/model"
)

type smallGroupUserService struct{}

var SmallGroupUserService = new(smallGroupUserService)

// ListByUserId 获取用户所加入的群组
func (*smallGroupUserService) ListByUserId(ctx context.Context, userId int64) ([]model.Group, error) {
	groups, err := dao.GroupUserDao.ListByUserId(userId)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// GetUsers 获取群组的所有用户信息
func (*smallGroupUserService) GetUsers(ctx context.Context, groupId int64) ([]model.GroupUser, error) {
	users, err := cache.GroupUserCache.Get(groupId)
	if err != nil {
		return nil, err
	}

	if users != nil {
		return users, nil
	}

	users, err = dao.GroupUserDao.ListUser(groupId)
	if err != nil {
		return nil, err
	}

	err = cache.GroupUserCache.Set(groupId, users)
	if err != nil {
		return nil, err
	}
	return users, err
}

// AddUser 给群组添加用户
func (*smallGroupUserService) AddUser(ctx context.Context, groupUser model.GroupUser) error {
	err := dao.GroupUserDao.Add(groupUser)
	if err != nil {
		return err
	}

	err = dao.GroupDao.UpdateUserNum(groupUser.GroupId, 1)
	if err != nil {
		return err
	}

	err = cache.GroupUserCache.Del(groupUser.GroupId)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser 从群组移除用户
func (*smallGroupUserService) DeleteUser(ctx context.Context, groupId, userId int64) error {
	err := dao.GroupUserDao.Delete(groupId, userId)
	if err != nil {
		return err
	}

	err = dao.GroupDao.UpdateUserNum(groupId, -1)
	if err != nil {
		return err
	}

	err = cache.GroupUserCache.Del(groupId)
	if err != nil {
		return err
	}

	return nil
}

// Update 更新群组用户信息
func (*smallGroupUserService) Update(ctx context.Context, user model.GroupUser) error {
	err := dao.GroupUserDao.Update(user)
	if err != nil {
		return err
	}

	err = cache.GroupUserCache.Del(user.GroupId)
	if err != nil {
		return err
	}
	return nil
}

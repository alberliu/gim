package service

import (
	"goim/logic/dao"
	"goim/logic/model"
	"goim/public/imctx"
	"goim/public/logger"
)

type friendService struct{}

var FriendService = new(friendService)

// List 获取用户好友列表
func (*friendService) ListUserFriend(ctx *imctx.Context, userId int64) ([]model.UserFriend, error) {
	friends, err := dao.FriendDao.ListUserFriend(ctx, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return friends, err
}

// Add 添加好友关系
func (*friendService) Add(ctx *imctx.Context, userId int64, friendId int64, label string) error {
	err := ctx.Session.Begin()
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	defer ctx.Session.Rollback()

	err = dao.FriendDao.Add(ctx, userId, friendId, label)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	err = dao.FriendDao.Add(ctx, friendId, userId, "")
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	ctx.Session.Commit()
	return nil
}

// Delete 删除好友关系
func (*friendService) Delete(ctx *imctx.Context, userId, friendId int64) error {
	err := ctx.Session.Begin()
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	defer ctx.Session.Rollback()

	err = dao.FriendDao.Delete(ctx, userId, friendId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	err = dao.FriendDao.Delete(ctx, friendId, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	ctx.Session.Commit()
	return nil
}

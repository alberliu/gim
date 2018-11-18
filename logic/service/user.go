package service

import (
	"database/sql"
	"goim/logic/dao"
	"goim/logic/model"
	"goim/public/imctx"
	"goim/public/imerror"
	"goim/public/logger"
)

type userService struct{}

var UserService = new(userService)

// Regist 注册
func (*userService) Regist(ctx *imctx.Context, deviceId int64, regist model.UserRegist) (*model.SignInResp, error) {
	err := ctx.Session.Begin()
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	defer ctx.Session.Rollback()

	// 添加用户
	user := model.User{
		Number:   regist.Number,
		Nickname: regist.Nickname,
		Sex:      regist.Sex,
		Avatar:   regist.Avatar,
		Password: regist.Password,
	}
	userId, err := dao.UserDao.Add(ctx, user)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	if userId == 0 {
		return nil, imerror.LErrNumberUsed
	}

	err = dao.UserSequenceDao.Add(ctx, userId, 0)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	err = dao.DeviceDao.UpdateUserId(ctx, deviceId, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	dao.DeviceSendSequenceDao.UpdateSendSequence(ctx, deviceId, 0)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	dao.DeviceSyncSequenceDao.UpdateSyncSequence(ctx, deviceId, 0)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	err = ctx.Session.Commit()
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	return &model.SignInResp{
		SendSequence: 0,
		SyncSequence: 0,
	}, nil
}

// SignIn 登录
func (*userService) SignIn(ctx *imctx.Context, deviceId int64, number string, password string) (*model.SignInResp, error) {
	err := ctx.Session.Begin()
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	defer ctx.Session.Rollback()
	// 设备验证

	// 用户验证
	user, err := dao.UserDao.GetByNumber(ctx, number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, imerror.LErrNameOrPassword
		}
		logger.Sugar.Error(err)
		return nil, err
	}
	if password != user.Password {
		return nil, imerror.LErrNameOrPassword
	}

	err = dao.DeviceDao.UpdateUserId(ctx, deviceId, user.Id)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	err = dao.DeviceSendSequenceDao.UpdateSendSequence(ctx, deviceId, 0)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	maxSyncSequence, err := dao.DeviceSyncSequenceDao.GetMaxSyncSequenceByUserId(ctx, user.Id)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	err = ctx.Session.Commit()
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return &model.SignInResp{
		SendSequence: 0,
		SyncSequence: maxSyncSequence,
	}, nil
}

// Get 获取用户信息
func (*userService) Get(ctx *imctx.Context, userId int64) (*model.User, error) {
	user, err := dao.UserDao.Get(ctx, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	user.Id = userId
	return user, err
}

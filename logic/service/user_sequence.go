package service

import (
	"goim/logic/dao"
	"goim/public/imctx"
	"goim/public/logger"
)

type userRequenceService struct{}

var UserRequenceService = new(userRequenceService)

// GetNext 获取下一个序列
func (*userRequenceService) GetNext(ctx *imctx.Context, userId int64) (int64, error) {
	err := ctx.Session.Begin()
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}
	ctx.Session.Rollback()

	err = dao.UserSequenceDao.Increase(ctx, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}

	sequence, err := dao.UserSequenceDao.GetSequence(ctx, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}

	return sequence, nil
}

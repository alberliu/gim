package dao

import (
	"goim/public/imctx"
	"goim/public/logger"
)

type userSequenceDao struct{}

var UserSequenceDao = new(userSequenceDao)

// Add 添加
func (*userSequenceDao) Add(ctx *imctx.Context, userId int64, sequence int64) error {
	_, err := ctx.Session.Exec("insert into t_user_sequence (user_id,sequence) values(?,?)", userId, sequence)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// Increase sequence++
func (*userSequenceDao) Increase(ctx *imctx.Context, userId int64) error {
	_, err := ctx.Session.Exec("update t_user_sequence set sequence = sequence + 1 where user_id = ?", userId)
	if err != nil {
		logger.Sugar.Error(err)
	}
	return err
}

// GetSequence 获取自增序列
func (*userSequenceDao) GetSequence(ctx *imctx.Context, userId int64) (int64, error) {
	var sequence int64
	err := ctx.Session.QueryRow("select sequence from t_user_sequence where user_id = ?", userId).
		Scan(&sequence)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err

	}
	return sequence, nil
}

package server

import (
	"database/sql"
	"goim/logic/dao"
	"goim/logic/db"
	"goim/logic/service"
	"goim/public/imctx"
	"goim/public/logger"
	"goim/public/transfer"
)

func Context() *imctx.Context {
	return imctx.NewContext(db.Factoty.GetSession())
}

type LogicRPCServer struct{}

// SignIn 处理设备登录
func (s *LogicRPCServer) SignIn(signIn transfer.SignIn, ack *transfer.SignInACK) error {
	ctx := Context()
	device, err := dao.DeviceDao.Get(ctx, signIn.DeviceId)
	if err == sql.ErrNoRows {
		ack.Code = transfer.CodeSignInFail
		ack.Message = "fail"
		return nil
	}

	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	if device.UserId == signIn.UserId && device.Token == signIn.Token {
		dao.DeviceDao.UpdateStatus(ctx, signIn.DeviceId, service.DeviceOnline)
		if err != nil {
			logger.Sugar.Error(err)
			return err
		}
		ack.Code = transfer.CodeSignInSuccess
		ack.Message = "success"
	} else {
		ack.Code = transfer.CodeSignInFail
		ack.Message = "fail"
	}

	logger.Sugar.Infow("设备登录",
		"device_id:", signIn.DeviceId,
		"user_id", signIn.UserId,
		"token", signIn.Token,
		"result", ack.Message)

	return nil
}

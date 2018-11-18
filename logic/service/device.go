package service

import (
	"goim/logic/dao"
	"goim/logic/model"

	"goim/public/imctx"
	"goim/public/logger"

	"github.com/satori/go.uuid"
)

const (
	DeviceOnline  = 1
	DeviceOffline = 0
)

type deviceService struct{}

var DeviceService = new(deviceService)

// Regist 注册设备
func (*deviceService) Regist(ctx *imctx.Context, device model.Device) (int64, string, error) {
	err := ctx.Session.Begin()
	if err != nil {
		logger.Sugar.Error(err)
		return 0, "", err
	}
	defer ctx.Session.Rollback()

	UUID := uuid.NewV4()
	device.Token = UUID.String()
	id, err := dao.DeviceDao.Add(ctx, device)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, "", err
	}

	err = dao.DeviceSendSequenceDao.Add(ctx, id, 0)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, "", err
	}

	err = dao.DeviceSyncSequenceDao.Add(ctx, id, 0)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, "", err
	}

	err = ctx.Session.Commit()
	if err != nil {
		logger.Sugar.Error(err)
		return 0, "", err
	}
	return id, device.Token, nil
}

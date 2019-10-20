package service

import (
	"goim/logic/cache"
	"goim/logic/dao"
	"goim/logic/model"
	"goim/public/imerror"
	"goim/public/util"

	"goim/public/imctx"
	"goim/public/logger"
)

const (
	DeviceOnline  = 1
	DeviceOffline = 0
)

type deviceService struct{}

var DeviceService = new(deviceService)

// Register 注册设备
func (*deviceService) Register(ctx *imctx.Context, device model.Device) (int64, error) {
	app, err := AppService.Get(ctx, device.AppId)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}

	if app == nil {
		return 0, imerror.ErrBadRequest
	}

	deviceId, err := util.DeviceIdUid.Get()
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}

	device.DeviceId = deviceId
	err = dao.DeviceDao.Add(ctx, device)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}

	err = dao.DeviceAckDao.Add(ctx, device.DeviceId, 0)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}

	return deviceId, nil
}

// ListOnlineByUserId 获取用户的所有在线设备
func (*deviceService) ListOnlineByUserId(ctx *imctx.Context, appId, userId int64) ([]model.Device, error) {
	devices, err := cache.UserDeviceCache.Get(appId, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	if devices != nil {
		return devices, nil
	}

	devices, err = dao.DeviceDao.ListOnlineByUserId(ctx, appId, userId)
	if err != nil {
		return nil, err
	}

	err = cache.UserDeviceCache.Set(appId, userId, devices)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	return devices, nil
}

// Online 设备上线
func (*deviceService) Online(ctx *imctx.Context, appId, deviceId, userId int64) error {
	err := dao.DeviceDao.UpdateUserIdAndStatus(ctx, deviceId, userId, DeviceOnline)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	err = cache.UserDeviceCache.Del(appId, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// Offline 设备离线
func (*deviceService) Offline(ctx *imctx.Context, appId, userId, deviceId int64) error {
	err := dao.DeviceDao.UpdateStatus(ctx, deviceId, DeviceOffline)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	err = cache.UserDeviceCache.Del(appId, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	err = cache.DeviceIPCache.Del(deviceId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

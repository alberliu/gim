package service

import (
	"context"
	"gim/internal/logic/cache"
	"gim/internal/logic/dao"
	"gim/internal/logic/model"
	"gim/pkg/gerrors"
	"gim/pkg/util"

	"gim/pkg/logger"
)

const (
	DeviceOnline  = 1
	DeviceOffline = 0
)

type deviceService struct{}

var DeviceService = new(deviceService)

// Register 注册设备
func (*deviceService) Register(ctx context.Context, device model.Device) (int64, error) {
	app, err := AppService.Get(ctx, device.AppId)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}

	if app == nil {
		return 0, gerrors.ErrBadRequest
	}

	deviceId, err := util.DeviceIdUid.Get()
	if err != nil {
		return 0, err
	}

	device.DeviceId = deviceId
	err = dao.DeviceDao.Add(device)
	if err != nil {
		return 0, err
	}

	err = dao.DeviceAckDao.Add(device.DeviceId, 0)
	if err != nil {
		return 0, err
	}

	return deviceId, nil
}

// ListOnlineByUserId 获取用户的所有在线设备
func (*deviceService) ListOnlineByUserId(ctx context.Context, appId, userId int64) ([]model.Device, error) {
	devices, err := cache.UserDeviceCache.Get(appId, userId)
	if err != nil {
		return nil, err
	}

	if devices != nil {
		return devices, nil
	}

	devices, err = dao.DeviceDao.ListOnlineByUserId(appId, userId)
	if err != nil {
		return nil, err
	}

	err = cache.UserDeviceCache.Set(appId, userId, devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// Online 设备上线
func (*deviceService) Online(ctx context.Context, appId, deviceId, userId int64, connectAddr string) error {
	err := dao.DeviceDao.UpdateUserIdAndStatus(deviceId, userId, DeviceOnline, connectAddr)
	if err != nil {
		return err
	}

	err = cache.UserDeviceCache.Del(appId, userId)
	if err != nil {
		return err
	}
	return nil
}

// Offline 设备离线
func (*deviceService) Offline(ctx context.Context, appId, userId, deviceId int64) error {
	err := dao.DeviceDao.UpdateStatus(deviceId, DeviceOffline)
	if err != nil {
		return err
	}

	err = cache.UserDeviceCache.Del(appId, userId)
	if err != nil {
		return err
	}

	err = cache.DeviceIPCache.Del(deviceId)
	if err != nil {
		return err
	}
	return nil
}

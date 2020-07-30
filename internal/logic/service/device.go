package service

import (
	"context"
	"gim/internal/logic/cache"
	"gim/internal/logic/dao"
	"gim/internal/logic/model"
)

const (
	DeviceOnline  = 1
	DeviceOffline = 0
)

type deviceService struct{}

var DeviceService = new(deviceService)

// Register 注册设备
func (*deviceService) Register(ctx context.Context, device model.Device) (int64, error) {
	id, err := dao.DeviceDao.Add(device)
	if err != nil {
		return 0, err
	}

	err = dao.DeviceAckDao.Add(id, 0)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (*deviceService) Get(ctx context.Context, deviceId int64) (*model.Device, error) {
	device, err := cache.DeviceCache.Get(deviceId)
	if err != nil {
		return nil, err
	}

	if device != nil {
		return device, nil
	}

	device, err = dao.DeviceDao.Get(deviceId)
	if err != nil {
		return nil, err
	}

	if device != nil {
		err = cache.DeviceCache.Set(device)
		if err != nil {
			return nil, err
		}
	}
	return device, nil
}

// ListOnlineByUserId 获取用户的所有在线设备
func (*deviceService) ListOnlineByUserId(ctx context.Context, userId int64) ([]model.Device, error) {
	devices, err := cache.UserDeviceCache.Get(userId)
	if err != nil {
		return nil, err
	}

	if devices != nil {
		return devices, nil
	}

	devices, err = dao.DeviceDao.ListOnlineByUserId(userId)
	if err != nil {
		return nil, err
	}

	err = cache.UserDeviceCache.Set(userId, devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// Online 设备上线
func (*deviceService) Online(ctx context.Context, deviceId, userId int64, connAddr string, connFd int64) error {
	err := dao.DeviceDao.UpdateUserIdAndStatus(deviceId, userId, DeviceOnline, connAddr, connFd)
	if err != nil {
		return err
	}

	err = cache.UserDeviceCache.Del(userId)
	if err != nil {
		return err
	}
	return nil
}

// Offline 设备离线
func (*deviceService) Offline(ctx context.Context, userId, deviceId int64) error {
	err := dao.DeviceDao.UpdateStatus(deviceId, DeviceOffline)
	if err != nil {
		return err
	}

	err = cache.UserDeviceCache.Del(userId)
	if err != nil {
		return err
	}
	return nil
}

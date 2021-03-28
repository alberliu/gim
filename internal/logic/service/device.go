package service

import (
	"context"
	"gim/internal/logic/cache"
	"gim/internal/logic/dao"
	"gim/internal/logic/model"
	"gim/pkg/logger"

	"go.uber.org/zap"
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
func (*deviceService) Online(ctx context.Context, deviceId, userId int64, connAddr string, clientAddr string) error {
	err := dao.DeviceDao.Update(deviceId, userId, DeviceOnline, connAddr, clientAddr)
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
func (*deviceService) Offline(ctx context.Context, userId, deviceId int64, clientAddr string) error {
	device, err := dao.DeviceDao.Get(deviceId)
	if err != nil {
		return err
	}
	if device == nil {
		logger.Logger.Warn("device is nil", zap.Int64("device_id", deviceId))
		return nil
	}
	if device.ClientAddr != clientAddr {
		return nil
	}

	err = dao.DeviceDao.UpdateStatus(deviceId, DeviceOffline)
	if err != nil {
		return err
	}

	err = cache.UserDeviceCache.Del(userId)
	if err != nil {
		return err
	}
	return nil
}

// ServerStop 设备离线
func (*deviceService) ServerStop(ctx context.Context, connAddr string) error {
	devices, err := dao.DeviceDao.ListOnlineByConnAddr(connAddr)
	if err != nil {
		return err
	}

	err = dao.DeviceDao.UpdateStatusByCoonAddr(connAddr, model.DeviceOffLine)
	if err != nil {
		return err
	}

	var userIds = make([]int64, 0, len(devices))
	for i := range devices {
		userIds = append(userIds, devices[i].UserId)
	}

	err = cache.UserDeviceCache.Del(userIds...)
	if err != nil {
		return err
	}
	return nil
}

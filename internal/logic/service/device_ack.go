package service

import (
	"context"
	"gim/internal/logic/dao"
)

type deviceAckService struct{}

var DeviceAckService = new(deviceAckService)

// Register 注册设备
func (*deviceAckService) Update(ctx context.Context, deviceId, ack int64) error {
	return dao.DeviceAckDao.Update(deviceId, ack)
}

func (*deviceAckService) GetMaxByUserId(ctx context.Context, userId int64) (int64, error) {
	return dao.DeviceAckDao.GetMaxByUserId(userId)
}

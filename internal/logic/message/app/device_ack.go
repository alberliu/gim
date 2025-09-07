package app

import (
	"context"

	"gim/internal/logic/message/repo"
)

var DeviceACKApp = new(deviceACKApp)

type deviceACKApp struct{}

// getMaxByUserId 根据用户id获取最大ack
func (*deviceACKApp) getMaxByUserId(ctx context.Context, userId uint64) (uint64, error) {
	acks, err := repo.DeviceACKRepo.Get(userId)
	if err != nil {
		return 0, err
	}

	var max uint64 = 0
	for i := range acks {
		if acks[i] > max {
			max = acks[i]
		}
	}
	return max, nil
}

// MessageAck 收到消息回执
func (*deviceACKApp) MessageAck(ctx context.Context, userId, deviceId, ack uint64) error {
	if ack <= 0 {
		return nil
	}
	return repo.DeviceACKRepo.Set(userId, deviceId, ack)
}

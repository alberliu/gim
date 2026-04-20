package domain

import "time"

type DeviceACK struct {
	DeviceID  uint64    // 设备ID
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
	UserID    uint64    // 用户ID
	ACK       uint64    // ack
}

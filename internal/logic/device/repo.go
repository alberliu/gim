package device

import (
	"errors"

	"gorm.io/gorm"

	"gim/internal/user/domain"
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

type repo struct{}

var Repo = new(repo)

// Get 获取设备
func (*repo) Get(deviceID uint64) (*Device, error) {
	var device Device
	err := db.DB.First(&device, "id = ?", deviceID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.ErrDeviceNotFound
	}
	return &device, err
}

// Save 保存设备信息
func (*repo) Save(device *Device) error {
	return db.DB.Save(&device).Error
}

// ListOnlineByUserId 获取用户的所有在线设备
func (*repo) ListOnlineByUserId(userIds []uint64) ([]Device, error) {
	var devices []Device
	err := db.DB.Find(&devices, "user_id in (?) and status = ?", userIds, OnLine).Error
	return devices, err
}

// ListOnlineByConnAddr 查询用户所有的在线设备
func (*repo) ListOnlineByConnAddr(connAddr string) ([]Device, error) {
	var devices []Device
	err := db.DB.Find(&devices, "conn_addr = ? and status = ?", connAddr, OnLine).Error
	return devices, err
}

// UpdateStatusOffline 更新设备为离线状态
func (*repo) UpdateStatusOffline(device Device) error {
	return db.DB.Model(&domain.Device{}).
		Where("id = ? and conn_addr = ?", device.ConnAddr, device.ConnAddr).
		Update("status", device.Status).Error
}

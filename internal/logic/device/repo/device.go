package repo

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"

	"gim/internal/logic/device/domain"
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

var DeviceRepo = new(deviceRepo)

type deviceRepo struct{}

// Get 获取设备
func (*deviceRepo) Get(deviceID uint64) (*domain.Device, error) {
	var device domain.Device
	err := db.DB.First(&device, "id = ?", deviceID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.ErrDeviceNotFound
	}
	return &device, err
}

// Save 保存设备信息
func (*deviceRepo) Save(device *domain.Device) error {
	return db.DB.Save(&device).Error
}

// ListByUserID 获取用户设备
func (r *deviceRepo) ListByUserID(userID uint64) ([]domain.Device, error) {
	var devices []domain.Device
	err := db.DB.Find(&devices, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	for i := range devices {
		isOnline, err := r.GetIsOnline(devices[i].ID)
		if err != nil {
			return nil, err
		}
		devices[i].IsOnline = isOnline
	}
	return devices, err
}

const deviceStatus = "deviceStatus:%d"

// SetOnline 设置在线
func (*deviceRepo) SetOnline(deviceID uint64) error {
	key := fmt.Sprintf(deviceStatus, deviceID)
	_, err := db.RedisCli.Set(key, "", 12*time.Minute).Result()
	return err
}

// SetOffline 设置在线
func (*deviceRepo) SetOffline(deviceID uint64) error {
	key := fmt.Sprintf(deviceStatus, deviceID)
	_, err := db.RedisCli.Del(key).Result()
	return err
}

func (*deviceRepo) GetIsOnline(deviceID uint64) (bool, error) {
	key := fmt.Sprintf(deviceStatus, deviceID)
	_, err := db.RedisCli.Get(key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

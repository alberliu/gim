package device

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"

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

// ListByUserID 获取用户设备
func (r *repo) ListByUserID(userID uint64) ([]Device, error) {
	var devices []Device
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
func (*repo) SetOnline(deviceID uint64) error {
	key := fmt.Sprintf(deviceStatus, deviceID)
	_, err := db.RedisCli.Set(key, "", 12*time.Minute).Result()
	return err
}

// SetOffline 设置在线
func (*repo) SetOffline(deviceID uint64) error {
	key := fmt.Sprintf(deviceStatus, deviceID)
	_, err := db.RedisCli.Del(key).Result()
	return err
}

func (*repo) GetIsOnline(deviceID uint64) (bool, error) {
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

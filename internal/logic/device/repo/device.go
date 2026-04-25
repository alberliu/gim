package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"gim/internal/logic/device/domain"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"gim/pkg/uredis"
)

const userDeviceKey = "userDevice:%d"

var DeviceRepo = new(deviceRepo)

type deviceRepo struct{}

// Get 获取设备
func (*deviceRepo) Get(ctx context.Context, deviceID uint64) (*domain.Device, error) {
	var device domain.Device
	err := db.DB.WithContext(ctx).First(&device, "id = ?", deviceID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.ErrDeviceNotFound
	}
	return &device, err
}

// Save 保存设备信息
func (*deviceRepo) Save(ctx context.Context, device *domain.Device) error {
	err := db.DB.WithContext(ctx).Save(&device).Error
	if err != nil {
		return err
	}
	key := fmt.Sprintf(userDeviceKey, device.UserID)
	return db.RedisCli.Del(ctx, key).Err()
}

// ListByUserID 获取用户设备
func (r *deviceRepo) ListByUserID(ctx context.Context, userID uint64) ([]domain.Device, error) {
	key := fmt.Sprintf(userDeviceKey, userID)
	devices, err := uredis.Get(db.RedisCli, ctx, key, 24*time.Hour, func() (*[]domain.Device, error) {
		devices, err := gorm.G[domain.Device](db.DB).Where("user_id = ?", userID).Find(ctx)
		if err != nil {
			return nil, err
		}
		return &devices, nil
	})
	if err != nil {
		return nil, err
	}

	// Status 是动态的，用 MGet 批量获取减少 Redis 往返
	if len(*devices) > 0 {
		keys := make([]string, len(*devices))
		for i := range *devices {
			keys[i] = fmt.Sprintf(deviceStatus, (*devices)[i].ID)
		}
		vals, err := db.RedisCli.MGet(ctx, keys...).Result()
		if err != nil {
			return nil, err
		}
		for i := range *devices {
			if vals[i] != nil {
				(*devices)[i].Status = domain.StatusOnline
			}
		}
	}
	return *devices, nil
}

const deviceStatus = "deviceStatus:%d"

// SetStatus 设置在线
func (*deviceRepo) SetStatus(ctx context.Context, deviceID uint64, status domain.Status) error {
	key := fmt.Sprintf(deviceStatus, deviceID)
	var err error
	if status == domain.StatusOnline {
		_, err = db.RedisCli.Set(ctx, key, "", 12*time.Minute).Result()
	} else {
		_, err = db.RedisCli.Del(ctx, key).Result()
	}
	return err
}

// GetStatus 获取状态
func (*deviceRepo) GetStatus(ctx context.Context, deviceID uint64) (domain.Status, error) {
	key := fmt.Sprintf(deviceStatus, deviceID)
	_, err := db.RedisCli.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return domain.StatusOffline, nil
	}
	if err != nil {
		return domain.StatusOffline, err
	}
	return domain.StatusOnline, nil
}

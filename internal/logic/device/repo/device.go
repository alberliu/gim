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

const deviceStatusTTL = 12 * time.Minute

// offlineScript 仅当当前 value 与传入 token 一致时才 DEL，避免旧连接关闭时
// 误把新连接已经写入的在线态抹掉。
var offlineScript = redis.NewScript(`
if redis.call('GET', KEYS[1]) == ARGV[1] then
    return redis.call('DEL', KEYS[1])
end
return 0
`)

// SetOnline 写入在线态，value 为本次连接的 token，覆盖式。
func (*deviceRepo) SetOnline(ctx context.Context, deviceID uint64, token string) error {
	key := fmt.Sprintf(deviceStatus, deviceID)
	return db.RedisCli.Set(ctx, key, token, deviceStatusTTL).Err()
}

// RefreshOnline 仅刷新 TTL，不改 value，供心跳调用，避免覆盖当前会话 token。
func (*deviceRepo) RefreshOnline(ctx context.Context, deviceID uint64) error {
	key := fmt.Sprintf(deviceStatus, deviceID)
	return db.RedisCli.Expire(ctx, key, deviceStatusTTL).Err()
}

// SetOffline 仅当当前 value 与 token 一致时才删除（fencing），
// 这样旧连接发来的 Offline 不会把新连接的在线态吃掉。
func (*deviceRepo) SetOffline(ctx context.Context, deviceID uint64, token string) error {
	key := fmt.Sprintf(deviceStatus, deviceID)
	_, err := offlineScript.Run(ctx, db.RedisCli, []string{key}, token).Result()
	if errors.Is(err, redis.Nil) {
		return nil
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

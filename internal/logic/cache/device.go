package cache

import (
	"gim/internal/logic/model"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	DeviceKey    = "device:"
	DeviceExpire = 2 * time.Hour
)

type deviceCache struct{}

var DeviceCache = new(deviceCache)

// Get 获取设备缓存
func (c *deviceCache) Get(deviceId int64) (*model.Device, error) {
	var device model.Device
	err := RedisUtil.Get(DeviceKey+strconv.FormatInt(deviceId, 10), &device)
	if err != nil && err != redis.Nil {
		return nil, gerrors.WrapError(err)
	}
	if err == redis.Nil {
		return nil, nil
	}
	return &device, nil
}

// Set 设置设备缓存
func (c *deviceCache) Set(device *model.Device) error {
	err := RedisUtil.Set(DeviceKey+strconv.FormatInt(device.Id, 10), device, DeviceExpire)
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Del 删除设备缓存
func (c *deviceCache) Del(deviceId int64) error {
	_, err := db.RedisCli.Del(DeviceKey + strconv.FormatInt(deviceId, 10)).Result()
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

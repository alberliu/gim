package cache

import (
	"gim/internal/logic/db"
	"gim/internal/logic/model"
	"gim/pkg/gerrors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	DeviceKey    = "user:device:"
	DeviceExpire = 2 * time.Hour
)

type userDeviceCache struct{}

var UserDeviceCache = new(userDeviceCache)

func (c *userDeviceCache) Key(appId, userId int64) string {
	return DeviceKey + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(userId, 10)
}

// Get 获取指定用户的所有在线设备
func (c *userDeviceCache) Get(appId, userId int64) ([]model.Device, error) {
	var devices []model.Device
	err := get(c.Key(appId, userId), &devices)
	if err != nil && err != redis.Nil {
		return nil, gerrors.WrapError(err)
	}

	if err == redis.Nil {
		return nil, nil
	}
	return devices, nil
}

// Set 将指定用户的所有在线设备存入缓存
func (c *userDeviceCache) Set(appId, userId int64, devices []model.Device) error {
	err := set(c.Key(appId, userId), devices, DeviceExpire)
	return gerrors.WrapError(err)
}

// Del 删除某一用户的在线设备列表
func (c *userDeviceCache) Del(appId, userId int64) error {
	_, err := db.RedisCli.Del(c.Key(appId, userId)).Result()
	return gerrors.WrapError(err)
}

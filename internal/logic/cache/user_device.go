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
	UserDeviceKey    = "user_device:"
	UserDeviceExpire = 2 * time.Hour
)

type userDeviceCache struct{}

var UserDeviceCache = new(userDeviceCache)

// Get 获取指定用户的所有在线设备
func (c *userDeviceCache) Get(userId int64) ([]model.Device, error) {
	var devices []model.Device
	err := RedisUtil.Get(UserDeviceKey+strconv.FormatInt(userId, 10), &devices)
	if err != nil && err != redis.Nil {
		return nil, gerrors.WrapError(err)
	}

	if err == redis.Nil {
		return nil, nil
	}
	return devices, nil
}

// Set 将指定用户的所有在线设备存入缓存
func (c *userDeviceCache) Set(userId int64, devices []model.Device) error {
	err := RedisUtil.Set(UserDeviceKey+strconv.FormatInt(userId, 10), devices, UserDeviceExpire)
	return gerrors.WrapError(err)
}

// Del 删除某一用户的在线设备列表
func (c *userDeviceCache) Del(userId int64) error {
	_, err := db.RedisCli.Del(UserDeviceKey + strconv.FormatInt(userId, 10)).Result()
	return gerrors.WrapError(err)
}

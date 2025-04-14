package device

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis"

	"gim/pkg/db"
	"gim/pkg/gerrors"
)

const (
	UserDeviceKey    = "user_device:"
	UserDeviceExpire = 2 * time.Hour
)

type userDeviceCache struct{}

var UserDeviceCache = new(userDeviceCache)

// Get 获取指定用户的所有在线设备
func (c *userDeviceCache) Get(userId int64) ([]Device, error) {
	var devices []Device
	err := db.RedisUtil.Get(UserDeviceKey+strconv.FormatInt(userId, 10), &devices)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, gerrors.WrapError(err)
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	return devices, nil
}

// Set 将指定用户的所有在线设备存入缓存
func (c *userDeviceCache) Set(userId int64, devices []Device) error {
	err := db.RedisUtil.Set(UserDeviceKey+strconv.FormatInt(userId, 10), devices, UserDeviceExpire)
	return gerrors.WrapError(err)
}

// Del 删除用户的在线设备列表
func (c *userDeviceCache) Del(userId int64) error {
	key := UserDeviceKey + strconv.FormatInt(userId, 10)
	_, err := db.RedisCli.Del(key).Result()
	return gerrors.WrapError(err)
}

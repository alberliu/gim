package cache

import (
	"gim/logic/db"
	"gim/public/logger"
	"strconv"

	"github.com/go-redis/redis"
)

const (
	DeviceIPKey = "device_ip:"
)

type deviceIPCache struct{}

var DeviceIPCache = new(deviceIPCache)

func (c *deviceIPCache) Key(deviceId int64) string {
	return DeviceIPKey + strconv.FormatInt(deviceId, 10)
}

// Get 获取设备所建立长连接的主机IP
func (c *deviceIPCache) Get(deviceId int64) (string, error) {
	ip, err := db.RedisCli.Get(DeviceIPKey + strconv.FormatInt(deviceId, 10)).Result()
	if err != nil && err != redis.Nil {
		logger.Sugar.Error(err)
		return "", err
	}
	if err == redis.Nil {
		return "", nil
	}
	return ip, nil
}

// Set 设置设备所建立长连接的主机IP
func (c *deviceIPCache) Set(deviceId int64, ip string) error {
	_, err := db.RedisCli.Set(DeviceIPKey+strconv.FormatInt(deviceId, 10), ip, 0).Result()
	if err != nil {
		logger.Sugar.Error(err)
	}
	return err
}

// Del 删除设备所建立长连接的主机IP
func (c *deviceIPCache) Del(deviceId int64) error {
	_, err := db.RedisCli.Del(DeviceIPKey + strconv.FormatInt(deviceId, 10)).Result()
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	return nil
}

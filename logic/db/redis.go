package db

import (
	"goim/conf"

	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

const DeviceIdPre = "connect:device_id:"

func init() {
	addr := conf.RedisIP

	RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		logs.Error("redis err ")
		panic(err)
	}
}

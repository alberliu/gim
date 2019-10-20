package db

import (
	"gim/conf"

	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
)

var RedisCli *redis.Client

func InitDB() {
	addr := conf.RedisIP

	RedisCli = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	_, err := RedisCli.Ping().Result()
	if err != nil {
		logs.Error("redis err ")
		panic(err)
	}
}

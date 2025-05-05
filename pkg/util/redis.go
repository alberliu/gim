package util

import (
	"time"

	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
)

type RedisUtil struct {
	client *redis.Client
}

func NewRedisUtil(client *redis.Client) *RedisUtil {
	return &RedisUtil{client: client}
}

// Set 将指定值设置到redis中，使用json的序列化方式
func (u *RedisUtil) Set(key string, value interface{}, duration time.Duration) error {
	bytes, err := jsoniter.Marshal(value)
	if err != nil {
		return err
	}

	return u.client.Set(key, bytes, duration).Err()
}

// Get 从redis中读取指定值，使用json的反序列化方式
func (u *RedisUtil) Get(key string, value interface{}) error {
	bytes, err := u.client.Get(key).Bytes()
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(bytes, value)
}

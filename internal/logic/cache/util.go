package cache

import (
	"gim/internal/logic/db"
	"gim/pkg/logger"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// set 将指定值设置到redis中，使用json的序列化方式
func set(key string, value interface{}, duration time.Duration) error {
	bytes, err := jsoniter.Marshal(value)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	err = db.RedisCli.Set(key, bytes, duration).Err()
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// get 从redis中读取指定值，使用json的反序列化方式
func get(key string, value interface{}) error {
	bytes, err := db.RedisCli.Get(key).Bytes()
	if err != nil {
		return err
	}
	err = jsoniter.Unmarshal(bytes, value)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

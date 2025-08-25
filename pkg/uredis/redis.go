package uredis

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/go-redis/redis"
)

type Client struct {
	*redis.Client
}

func NewClient(addr, password string) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       0,
		Password: password,
	})

	_, err := client.Ping().Result()
	if err != nil {
		slog.Error("redis ping error", "error", err)
		panic(err)
	}
	return &Client{Client: client}
}

// SetObject 将指定值设置到redis中，使用json的序列化方式
func (c *Client) SetObject(key string, value any, duration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.Set(key, bytes, duration).Err()
}

// GetObject 从redis中读取指定值，使用json的反序列化方式
func (c *Client) GetObject(key string, value any) error {
	bytes, err := c.Get(key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, value)
}

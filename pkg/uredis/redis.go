package uredis

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
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

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		slog.Error("redis ping error", "error", err)
		panic(err)
	}
	return &Client{Client: client}
}

// SetAny 将指定值设置到redis中，使用json的序列化方式
func (c *Client) SetAny(ctx context.Context, key string, value any, duration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.Set(ctx, key, bytes, duration).Err()
}

// GetAny 从redis中读取指定值，使用json的反序列化方式
func (c *Client) GetAny(ctx context.Context, key string, value any) error {
	bytes, err := c.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, value)
}

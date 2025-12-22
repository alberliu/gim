package uredis

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

// group 用于防止缓存击穿，合并同一进程内的并发请求
var group singleflight.Group

func Get[T any](client *Client, ctx context.Context, key string, ttl time.Duration, fallback func() (*T, error)) (*T, error) {
	t := new(T)
	err := client.GetAny(ctx, key, t)
	if err == nil {
		return t, nil
	}
	if !errors.Is(err, redis.Nil) {
		slog.ErrorContext(ctx, "Get", "error", err, "key", key)
	}

	result, err, _ := group.Do(key, func() (any, error) {
		t, err := fallback()
		if err != nil {
			return nil, err
		}
		err = client.SetAny(ctx, key, t, ttl)
		if err != nil {
			slog.ErrorContext(ctx, "Get", "error", err, "key", key)
		}
		return t, nil
	})

	if err != nil {
		return nil, err
	}
	return result.(*T), nil
}

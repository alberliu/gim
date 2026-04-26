package uredis

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"maps"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"

	"gim/pkg/uslices"
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
		slog.ErrorContext(ctx, "cache get failed", "error", err, "key", key)
	}

	result, err, _ := group.Do(key, func() (any, error) {
		t, err := fallback()
		if err != nil {
			return nil, err
		}
		err = client.SetAny(ctx, key, t, ttl)
		if err != nil {
			slog.ErrorContext(ctx, "cache set failed", "error", err, "key", key)
		}
		return t, nil
	})

	if err != nil {
		return nil, err
	}
	return result.(*T), nil
}

func MGet[K comparable, V any](client *Client, ctx context.Context, ids []K, getKey func(K) string,
	ttl time.Duration, fallback func(ids []K) (map[K]V, error)) (map[K]V, error) {
	if len(ids) == 0 {
		return make(map[K]V), nil
	}

	ids = uslices.Unique(ids)
	hash := make(map[K]V, len(ids))
	var miss []K

	keys := make([]string, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, getKey(id))
	}
	result, err := client.MGet(ctx, keys...).Result()
	if err != nil {
		miss = ids
		slog.ErrorContext(ctx, "redis mget failed", "error", err, "keys", keys)
	} else {
		for i, id := range ids {
			if result[i] == nil {
				miss = append(miss, id)
			} else {
				var v V
				err := json.Unmarshal([]byte(result[i].(string)), &v)
				if err != nil {
					slog.ErrorContext(ctx, "unmarshal cache value failed", "error", err, "key", keys[i])
					miss = append(miss, id)
				} else {
					hash[id] = v
				}
			}
		}
	}
	if len(miss) == 0 {
		return hash, nil
	}

	missHash, err := fallback(miss)
	if err != nil {
		return nil, err
	}
	maps.Copy(hash, missHash)

	pipe := client.Pipeline()
	for k, v := range missHash {
		key := getKey(k)
		buf, err := json.Marshal(v)
		if err != nil {
			slog.ErrorContext(ctx, "marshal cache value failed", "error", err, "key", key)
			continue
		}
		pipe.Set(ctx, key, buf, ttl)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "pipeline exec failed", "error", err)
	}
	return hash, nil
}

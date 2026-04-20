package repo

import (
	"context"
	"fmt"

	"gim/pkg/db"
)

const (
	SeqObjectTypeUser = 1 // 用户
)

type seqRepo struct{}

var SeqRepo = new(seqRepo)

// Incr 自增seq,并且获取自增后的值
func (*seqRepo) Incr(ctx context.Context, userID uint64) (uint64, error) {
	key := fmt.Sprintf("seq:%d", userID)
	val, err := db.RedisCli.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return uint64(val), nil
}

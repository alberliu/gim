package service

import (
	"context"
	"gim/internal/logic/cache"
)

type seqService struct{}

var SeqService = new(seqService)

// GetUserNext 获取下一个序列号
func (*seqService) GetUserNext(ctx context.Context, appId, userId int64) (int64, error) {
	return cache.SeqCache.Incr(cache.SeqCache.UserKey(appId, userId))
}

// GetGroupNext 获取下一个序列号
func (*seqService) GetGroupNext(ctx context.Context, appId, groupId int64) (int64, error) {
	return cache.SeqCache.Incr(cache.SeqCache.UserKey(appId, groupId))
}

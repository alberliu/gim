package service

import (
	"gim/logic/cache"
	"gim/public/imctx"
)

type seqService struct{}

var SeqService = new(seqService)

// GetUserNext 获取下一个序列号
func (*seqService) GetUserNext(ctx *imctx.Context, appId, userId int64) (int64, error) {
	return cache.SeqCache.Incr(cache.SeqCache.UserKey(appId, userId))
}

// GetGroupNext 获取下一个序列号
func (*seqService) GetGroupNext(ctx *imctx.Context, appId, groupId int64) (int64, error) {
	return cache.SeqCache.Incr(cache.SeqCache.UserKey(appId, groupId))
}

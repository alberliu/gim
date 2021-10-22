package service

import (
	"context"
	"gim/internal/logic/domain/message/repo"
)

type seqService struct{}

var SeqService = new(seqService)

// GetUserNext 获取下一个序列号
func (*seqService) GetUserNext(ctx context.Context, userId int64) (int64, error) {
	return repo.SeqRepo.Incr(repo.SeqObjectTypeUser, userId)
}

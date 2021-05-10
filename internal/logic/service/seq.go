package service

import (
	"context"
	"gim/internal/logic/dao"
)

type seqService struct{}

var SeqService = new(seqService)

// GetUserNext 获取下一个序列号
func (*seqService) GetUserNext(ctx context.Context, userId int64) (int64, error) {
	return dao.SeqDao.Incr(dao.SeqObjectTypeUser, userId)
}

// GetRoomNext 获取下一个序列号
func (*seqService) GetRoomNext(ctx context.Context, roomId int64) (int64, error) {
	return dao.SeqDao.Incr(dao.SeqObjectTypeRoom, roomId)
}

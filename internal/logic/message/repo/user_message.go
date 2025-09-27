package repo

import (
	"context"

	"gim/internal/logic/message/domain"
	"gim/pkg/db"
)

var UserMessageRepo = new(userMessageRepo)

type userMessageRepo struct{}

// Create 创建
func (d *userMessageRepo) Create(ctx context.Context, message *domain.UserMessage) error {
	seq, err := SeqRepo.Incr(ctx, SeqObjectTypeUser, message.UserID)
	if err != nil {
		return err
	}
	message.Seq = seq
	return db.DB.WithContext(ctx).Create(&message).Error
}

// ListBySeq 根据类型和id查询大于序号大于seq的消息
func (d *userMessageRepo) ListBySeq(ctx context.Context, userId, seq uint64, limit int64) ([]domain.UserMessage, bool, error) {
	var messages []domain.UserMessage
	limitint := int(limit)
	err := db.DB.WithContext(ctx).Table("user_message").
		Where("user_id = ? and seq > ?", userId, seq).
		Limit(limitint + 1).
		Preload("Message").
		Find(&messages).Error
	if err != nil {
		return nil, false, err
	}
	hasMore := len(messages) > limitint
	if hasMore {
		messages = messages[:limitint]
	}
	return messages, hasMore, nil
}

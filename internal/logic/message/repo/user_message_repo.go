package repo

import (
	"gim/internal/logic/message/domain"
	"gim/pkg/db"
)

var UserMessageRepo = new(userMessageRepo)

type userMessageRepo struct{}

// Create 创建
func (d *userMessageRepo) Create(message *domain.UserMessage) error {
	return db.DB.Create(&message).Error
}

// ListBySeq 根据类型和id查询大于序号大于seq的消息
func (d *userMessageRepo) ListBySeq(userId, seq uint64, limit int64) ([]domain.UserMessage, bool, error) {
	DB := db.DB.Table("user_message").
		Where("user_id = ? and seq > ?", userId, seq)

	var count int64
	err := DB.Count(&count).Error
	if err != nil {
		return nil, false, err
	}
	if count == 0 {
		return nil, false, nil
	}

	var messages []domain.UserMessage
	err = DB.Limit(int(limit)).Preload("Message").Find(&messages).Error
	if err != nil {
		return nil, false, err
	}
	return messages, count > limit, nil
}

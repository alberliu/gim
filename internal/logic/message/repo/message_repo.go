package repo

import (
	"gim/internal/logic/message/domain"
	"gim/pkg/db"
)

var MessageRepo = new(messageRepo)

type messageRepo struct{}

func (*messageRepo) Save(message *domain.Message) error {
	return db.DB.Create(&message).Error
}

func (*messageRepo) GetByIDs(ids []int64) ([]domain.Message, error) {
	var messages []domain.Message
	err := db.DB.Find(&messages, "id in (?)", ids).Error
	return messages, err
}

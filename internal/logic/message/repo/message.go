package repo

import (
	"context"

	"gim/internal/logic/message/domain"
	"gim/pkg/db"
)

var MessageRepo = new(messageRepo)

type messageRepo struct{}

func (*messageRepo) Save(ctx context.Context, message *domain.Message) error {
	return db.DB.WithContext(ctx).Create(&message).Error
}

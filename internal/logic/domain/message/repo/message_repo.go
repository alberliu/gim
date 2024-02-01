package repo

import (
	"fmt"
	"gim/internal/logic/domain/message/model"
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

type messageRepo struct{}

var MessageRepo = new(messageRepo)

func (*messageRepo) tableName(userId int64) string {
	return fmt.Sprintf("message")
}

// Save 插入一条消息
func (d *messageRepo) Save(message model.Message) error {
	err := db.DB.Table(d.tableName(message.UserId)).Create(&message).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// ListBySeq 根据类型和id查询大于序号大于seq的消息
func (d *messageRepo) ListBySeq(userId, seq, limit int64) ([]model.Message, bool, error) {
	DB := db.DB.Table(d.tableName(userId)).
		Where("user_id = ? and seq > ?", userId, seq)

	var count int64
	err := DB.Count(&count).Error
	if err != nil {
		return nil, false, gerrors.WrapError(err)
	}
	if count == 0 {
		return nil, false, nil
	}

	var messages []model.Message
	err = DB.Limit(limit).Find(&messages).Error
	if err != nil {
		return nil, false, gerrors.WrapError(err)
	}
	return messages, count > limit, nil
}

package dao

import (
	"fmt"
	"gim/internal/logic/model"
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

const messageTableNum = 1

type messageDao struct{}

var MessageDao = new(messageDao)

func (*messageDao) tableName(objectId int64) string {
	return fmt.Sprintf("message_%03d", objectId%messageTableNum)
}

// Add 插入一条消息
func (d *messageDao) Add(message model.Message) error {
	err := db.DB.Table(d.tableName(message.ObjectId)).Create(&message).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// ListBySeq 根据类型和id查询大于序号大于seq的消息
func (d *messageDao) ListBySeq(objectType, objectId, seq, limit int64) ([]model.Message, bool, error) {
	db := db.DB.Table(d.tableName(objectId)).
		Where("object_type = ? and object_id = ? and seq > ?", objectType, objectId, seq)

	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, false, gerrors.WrapError(err)
	}
	if count == 0 {
		return nil, false, nil
	}

	var messages []model.Message
	err = db.Limit(limit).Find(&messages).Error
	if err != nil {
		return nil, false, gerrors.WrapError(err)
	}
	return messages, count > limit, nil
}

package dao

import (
	"gim/internal/logic/model"
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

type messageDao struct{}

var MessageDao = new(messageDao)

// Add 插入一条消息
func (*messageDao) Add(tableName string, message model.Message) error {
	err := db.DB.Table(tableName).Create(&message).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// ListBySeq 根据类型和id查询大于序号大于seq的消息
func (*messageDao) ListBySeq(tableName string, objectType, objectId, seq int64) ([]model.Message, error) {
	var messages []model.Message
	err := db.DB.Table(tableName).Find(&messages, "object_type = ? and object_id = ? and seq > ?", objectType, objectId, seq).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return messages, nil
}

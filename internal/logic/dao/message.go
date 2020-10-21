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
func (*messageDao) ListBySeq(tableName string, objectType, objectId, seq, limit int64) ([]model.Message, bool, error) {
	db := db.DB.Table(tableName).
		Where("object_type = ? and object_id = ? and seq > ?", objectType, objectId, seq)

	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, false, gerrors.WrapError(err)
	}

	var messages []model.Message
	err = db.Limit(limit).Find(&messages).Error
	if err != nil {
		return nil, false, gerrors.WrapError(err)
	}
	return messages, count > limit, nil
}

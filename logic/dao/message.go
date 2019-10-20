package dao

import (
	"fmt"
	"gim/logic/db"
	"gim/logic/model"
	"gim/public/imctx"
	"gim/public/logger"
)

type messageDao struct{}

var MessageDao = new(messageDao)

// Add 插入一条消息
func (*messageDao) Add(ctx *imctx.Context, tableName string, message model.Message) error {
	sql := fmt.Sprintf(`insert into %s(message_id,app_id,object_type,object_id,sender_type,sender_id,sender_device_id,receiver_type,receiver_id,to_user_ids,type,content,seq,send_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, tableName)
	_, err := db.DBCli.Exec(sql, message.MessageId, message.AppId, message.ObjectType, message.ObjectId, message.SenderType, message.SenderId, message.SenderDeviceId,
		message.ReceiverType, message.ReceiverId, message.ToUserIds, message.Type, message.Content, message.Seq, message.SendTime)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// ListBySeq 根据类型和id查询大于序号大于seq的消息
func (*messageDao) ListBySeq(ctx *imctx.Context, tableName string, appId, objectType, objectId, seq int64) ([]model.Message, error) {
	sql := fmt.Sprintf(`select message_id,app_id,object_type,object_id,sender_type,sender_id,sender_device_id,receiver_type,receiver_id,to_user_ids,type,content,seq,send_time from %s where app_id = ? and object_type = ? and object_id = ? and seq > ?`, tableName)
	rows, err := db.DBCli.Query(sql, appId, objectType, objectId, seq)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	messages := make([]model.Message, 0, 5)
	for rows.Next() {
		message := new(model.Message)
		err := rows.Scan(&message.MessageId, &message.AppId, &message.ObjectType, &message.ObjectId, &message.SenderType, &message.SenderId, &message.SenderDeviceId, &message.ReceiverType,
			&message.ReceiverId, &message.ToUserIds, &message.Type, message.Content, &message.Seq, &message.SendTime)
		if err != nil {
			logger.Sugar.Error(err)
			return nil, err
		}
		messages = append(messages, *message)
	}
	return messages, nil
}

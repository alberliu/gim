package dao

import (
	"fmt"
	"gim/internal/logic/db"
	"gim/internal/logic/model"
	"gim/pkg/gerrors"
)

type messageDao struct{}

var MessageDao = new(messageDao)

// Add 插入一条消息
func (*messageDao) Add(tableName string, message model.Message) error {
	sql := fmt.Sprintf(`insert into %s(app_id,object_type,object_id,request_id,sender_type,sender_id,sender_device_id,receiver_type,receiver_id,to_user_ids,type,content,seq,send_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, tableName)
	_, err := db.DBCli.Exec(sql, message.AppId, message.ObjectType, message.ObjectId, message.RequestId, message.SenderType, message.SenderId, message.SenderDeviceId,
		message.ReceiverType, message.ReceiverId, message.ToUserIds, message.Type, message.Content, message.Seq, message.SendTime)
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// ListBySeq 根据类型和id查询大于序号大于seq的消息
func (*messageDao) ListBySeq(tableName string, appId, objectType, objectId, seq int64) ([]model.Message, error) {
	sql := fmt.Sprintf(`select app_id,object_type,object_id,request_id,sender_type,sender_id,sender_device_id,receiver_type,receiver_id,to_user_ids,type,content,seq,send_time from %s where app_id = ? and object_type = ? and object_id = ? and seq > ?`, tableName)
	rows, err := db.DBCli.Query(sql, appId, objectType, objectId, seq)
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	messages := make([]model.Message, 0, 5)
	for rows.Next() {
		message := new(model.Message)
		err := rows.Scan(&message.AppId, &message.ObjectType, &message.ObjectId, &message.RequestId, &message.SenderType, &message.SenderId, &message.SenderDeviceId, &message.ReceiverType,
			&message.ReceiverId, &message.ToUserIds, &message.Type, &message.Content, &message.Seq, &message.SendTime)
		if err != nil {
			return nil, gerrors.WrapError(err)
		}
		messages = append(messages, *message)
	}
	return messages, nil
}

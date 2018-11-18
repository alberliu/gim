package dao

import (
	"goim/logic/model"
	"goim/public/imctx"
	"goim/public/logger"
)

type messageDao struct{}

var MessageDao = new(messageDao)

// Add 插入一条消息
func (*messageDao) Add(ctx *imctx.Context, message model.Message) error {
	_, err := ctx.Session.Exec("insert into t_message(message_id,user_id,sender_type,sender_id,sender_device_id,receiver_type,receiver_id,type,content,send_time,sequence) values(?,?,?,?,?,?,?,?,?,?,?)",
		message.MessageId, message.UserId, message.SenderType, message.SenderId, message.SenderDeviceId,
		message.ReceiverType, message.ReceiverId, message.Type, message.Content, message.SendTime, message.Sequence)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// ListByUserIdAndSequence 根据用户id查询大于序号大于sequence的消息
func (*messageDao) ListByUserIdAndSequence(ctx *imctx.Context, userId int64, sequence int64) ([]*model.Message, error) {
	rows, err := ctx.Session.Query("select id,message_id,user_id,sender_type,sender_id,sender_device_id,receiver_type,receiver_id,type,content,sequence,send_time,create_time from t_message where user_id = ? and sequence >= ?",
		userId, sequence)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	messages := make([]*model.Message, 0, 5)
	for rows.Next() {
		message := new(model.Message)
		err := rows.Scan(&message.Id, &message.MessageId, &message.UserId, &message.SenderType, &message.SenderId, &message.SenderDeviceId, &message.ReceiverType,
			&message.ReceiverId, &message.Type, &message.Content, &message.Sequence, &message.SendTime, &message.CreateTime)
		if err != nil {
			logger.Sugar.Error(err)
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

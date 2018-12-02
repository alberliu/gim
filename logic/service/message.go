package service

import (
	"database/sql"
	"goim/logic/dao"
	"goim/logic/model"
	"goim/logic/rpc/connect_rpc"
	"goim/public/imctx"
	"goim/public/imerror"
	"goim/public/logger"
	"goim/public/transfer"
)

const (
	ReceiverUser  = 1 // 接收者类型为用户
	ReceiverGroup = 2 // 接收者类型为群组
)

const (
	SenderTypeUser  = 1 // 用户发送
	SenderTypeOther = 2 // 其他发送，业务推送
)

type messageService struct{}

var MessageService = new(messageService)

// Add 添加消息
func (*messageService) Add(ctx *imctx.Context, message model.Message) error {
	return dao.MessageDao.Add(ctx, message)
}

// ListByUserIdAndSequence 查询消息
func (*messageService) ListByUserIdAndSequence(ctx *imctx.Context, userId int64, sequence int64) ([]*model.Message, error) {
	return dao.MessageDao.ListByUserIdAndSequence(ctx, userId, sequence)
}

// SendToUser 消息发送至用户
func (*messageService) SendToFriend(ctx *imctx.Context, send transfer.MessageSend) error {
	_, err := dao.FriendDao.Get(ctx, send.SenderUserId, send.ReceiverId)
	if err == sql.ErrNoRows {
		logger.Sugar.Error(ctx, send.SenderUserId, send.ReceiverId, "不是好友关系")
		return imerror.CErrNotIsFriend
	}
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	selfSequence, err := UserRequenceService.GetNext(ctx, send.SenderUserId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	selfMessage := model.Message{
		MessageId:      send.MessageId,
		UserId:         send.SenderUserId,
		SenderType:     SenderTypeUser,
		SenderId:       send.SenderUserId,
		SenderDeviceId: send.SenderDeviceId,
		ReceiverType:   int(send.ReceiverType),
		ReceiverId:     send.ReceiverId,
		Type:           int(send.Type),
		Content:        send.Content,
		Sequence:       selfSequence,
		SendTime:       send.SendTime,
	}

	// 发给发送者
	err = MessageService.SendToUser(ctx, send.SenderUserId, &selfMessage)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	friendSequence, err := UserRequenceService.GetNext(ctx, send.ReceiverId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	friendMessage := model.Message{
		MessageId:      send.MessageId,
		UserId:         send.ReceiverId,
		SenderType:     SenderTypeUser,
		SenderId:       send.SenderUserId,
		SenderDeviceId: send.SenderDeviceId,
		ReceiverType:   int(send.ReceiverType),
		ReceiverId:     send.ReceiverId,
		Type:           int(send.Type),
		Content:        send.Content,
		Sequence:       friendSequence,
		SendTime:       send.SendTime,
	}
	// 发给接收者
	err = MessageService.SendToUser(ctx, send.ReceiverId, &friendMessage)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// SendToGroup 消息发送至群组
func (*messageService) SendToGroup(ctx *imctx.Context, send transfer.MessageSend) error {
	in, err := dao.GroupUserDao.UserInGroup(ctx, send.ReceiverId, send.SenderUserId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	if !in {
		logger.Sugar.Error(ctx, send.SenderUserId, send.ReceiverId, "不在群组内")
		return imerror.CErrNotInGroup
	}

	group, err := GroupService.Get(ctx, send.ReceiverId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	// 持久化到数据库
	for _, user := range group.GroupUser {
		sequence, err := UserRequenceService.GetNext(ctx, user.UserId)
		if err != nil {
			logger.Sugar.Error(err)
			return err
		}
		message := model.Message{
			MessageId:      send.MessageId,
			UserId:         user.UserId,
			SenderType:     SenderTypeUser,
			SenderId:       send.SenderUserId,
			SenderDeviceId: send.SenderDeviceId,
			ReceiverType:   int(send.ReceiverType),
			ReceiverId:     send.ReceiverId,
			Type:           int(send.Type),
			Content:        send.Content,
			Sequence:       sequence,
			SendTime:       send.SendTime,
		}

		err = MessageService.SendToUser(ctx, user.UserId, &message)
		if err != nil {
			logger.Sugar.Error(err)
			return err
		}

	}
	return nil
}

// SendToUser 消息发送至用户
func (*messageService) SendToUser(ctx *imctx.Context, userId int64, message *model.Message) error {
	err := MessageService.Add(ctx, *message)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	messageItem := transfer.MessageItem{
		MessageId:      message.MessageId,
		SenderType:     message.SenderType,
		SenderId:       message.SenderId,
		SenderDeviceId: message.SenderDeviceId,
		ReceiverType:   message.ReceiverType,
		ReceiverId:     message.ReceiverId,
		Type:           message.Type,
		Content:        message.Content,
		Sequence:       message.Sequence,
		SendTime:       message.SendTime,
	}

	// 查询用户在线设备
	devices, err := dao.DeviceDao.ListOnlineByUserId(ctx, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	for _, v := range devices {
		message := transfer.Message{DeviceId: v.Id, Type: transfer.MessageTypeMail, Messages: []transfer.MessageItem{messageItem}}
		connect_rpc.ConnectRPC.SendMessage(message)

		logger.Sugar.Infow("消息投递",
			"device_id:", message.DeviceId,
			"user_id", userId,
			"type", message.Type,
			"messages", message.GetLog())
	}
	return nil
}

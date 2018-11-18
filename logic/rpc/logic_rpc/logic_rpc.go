package logic_rpc

import (
	"database/sql"
	"goim/logic/dao"
	"goim/logic/rpc/connect_rpc"
	"goim/logic/service"
	"goim/public/imctx"
	"goim/public/imerror"
	"goim/public/lib"
	"goim/public/logger"
	"goim/public/transfer"
)

type logicRPC struct{}

var LogicRPC = new(logicRPC)

// SignIn 处理设备登录
func (s *logicRPC) SignIn(ctx *imctx.Context, signIn transfer.SignIn) (*transfer.SignInACK, error) {
	device, err := dao.DeviceDao.Get(ctx, signIn.DeviceId)
	if err == sql.ErrNoRows {
		return &transfer.SignInACK{
			Code:    transfer.CodeSignInFail,
			Message: "fail",
		}, nil
	}

	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	var code int
	var message string
	if device.UserId == signIn.UserId && device.Token == signIn.Token {
		dao.DeviceDao.UpdateStatus(ctx, signIn.DeviceId, service.DeviceOnline)
		if err != nil {
			logger.Sugar.Error(err)
			return nil, err
		}
		code = transfer.CodeSignInSuccess
		message = "success"
	} else {
		code = transfer.CodeSignInFail
		message = "fail"
	}

	logger.Sugar.Infow("设备登录",
		"device_id:", signIn.DeviceId,
		"user_id", signIn.UserId,
		"token", signIn.Token,
		"result", message)

	return &transfer.SignInACK{
		Code:    code,
		Message: message,
	}, err
}

// SyncTrigger 处理消息同步触发
func (s *logicRPC) SyncTrigger(ctx *imctx.Context, trigger transfer.SyncTrigger) error {
	logger.Sugar.Infow("同步触发",
		"device_id:", trigger.DeviceId,
		"user_id", trigger.UserId,
		"sync_sequence", trigger.SyncSequence)

	dbMessages, err := dao.MessageDao.ListByUserIdAndSequence(ctx, trigger.UserId, trigger.SyncSequence)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	messages := make([]transfer.MessageItem, 0, len(dbMessages))
	for _, v := range dbMessages {
		item := transfer.MessageItem{}

		item.MessageId = v.MessageId
		item.SenderType = v.SenderType
		item.SenderId = v.SenderId
		item.SenderDeviceId = v.SenderDeviceId
		item.ReceiverType = v.ReceiverType
		item.ReceiverId = v.ReceiverId
		item.Type = v.Type
		item.Content = v.Content
		item.Sequence = v.Sequence
		item.SendTime = v.SendTime

		messages = append(messages, item)
	}

	message := transfer.Message{DeviceId: trigger.DeviceId, Messages: messages}
	connect_rpc.ConnectRPC.SendMessage(message)

	logger.Sugar.Infow("消息同步",
		"device_id:", trigger.DeviceId,
		"user_id", trigger.UserId,
		"messages", message.GetLog())
	return nil
}

// MessageSend 处理消息发送
func (s *logicRPC) MessageSend(ctx *imctx.Context, send transfer.MessageSend) error {
	var err error
	send.MessageId = lib.Lid.Get()

	logger.Sugar.Infow("消息发送",
		"device_id", send.SenderDeviceId,
		"user_id", send.SenderUserId,
		"message_id", send.MessageId,
		"send_sequence", send.SendSequence)

	// 检查消息是否重复发送
	sendSequence, err := dao.DeviceSendSequenceDao.Get(ctx, send.SenderDeviceId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	ack := transfer.MessageSendACK{
		MessageId:    send.MessageId,
		DeviceId:     send.SenderDeviceId,
		SendSequence: send.SendSequence,
		Code:         imerror.CCodeSuccess,
	}
	if send.SendSequence <= sendSequence {
		// 消息发送回执
		err = connect_rpc.ConnectRPC.SendMessageSendACK(ack)
		if err != nil {
			logger.Sugar.Error(err)
		}
		return nil
	}
	err = dao.DeviceSendSequenceDao.UpdateSendSequence(ctx, send.SenderDeviceId, send.SendSequence)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	if send.ReceiverType == service.ReceiverUser {
		err = service.MessageService.SendToFriend(ctx, send)
	}
	if send.ReceiverType == service.ReceiverGroup {
		err = service.MessageService.SendToGroup(ctx, send)
	}

	if cerror, ok := err.(*imerror.CError); ok {
		ack.Code = cerror.Code
	} else {
		ack.Code = imerror.CErrUnkonw.Code
	}
	// 消息发送回执
	err = connect_rpc.ConnectRPC.SendMessageSendACK(ack)
	if err != nil {
		logger.Sugar.Error(err)
	}

	logger.Sugar.Infow("消息发送回执",
		"device_id", ack.DeviceId,
		"user_id", send.SenderUserId,
		"message_id", send.MessageId,
		"send_sequence", ack.SendSequence,
		"code", ack.Code,
	)

	return nil
}

// MessageACK 处理消息回执
func (s *logicRPC) MessageACK(ctx *imctx.Context, ack transfer.MessageACK) error {
	err := dao.DeviceSyncSequenceDao.UpdateSyncSequence(ctx, ack.DeviceId, ack.SyncSequence)
	if err != nil {
		logger.Sugar.Error(err)
	}

	logger.Sugar.Infow("消息投递回执",
		"device_id", ack.DeviceId,
		"user_id", ack.UserId,
		"message_id", ack.MessageId,
		"sync_sequence", ack.SyncSequence)

	return nil
}

// OffLine 处理设备离线
func (s *logicRPC) OffLine(ctx *imctx.Context, deviceId int64, userId int64) error {
	err := dao.DeviceDao.UpdateStatus(ctx, deviceId, service.DeviceOffline)
	if err != nil {
		logger.Sugar.Error(err)
	}

	logger.Sugar.Infow("设备离线", "device_id", deviceId, "user_id", userId)
	return nil
}

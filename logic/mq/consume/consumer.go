package consume

import (
	"time"

	"goim/logic/dao"
	"goim/logic/db"
	"goim/logic/service"
	"goim/public/imctx"
	"goim/public/imerror"
	"goim/public/lib"
	"goim/public/logger"
	"goim/public/transfer"

	"goim/conf"

	"goim/logic/mq/produce"

	"github.com/json-iterator/go"
	"github.com/nsqio/go-nsq"
)

// NsqConsumer 消费消息
func NsqConsumer(topic, channel string, handle func(message *nsq.Message) error, concurrency int) {
	config := nsq.NewConfig()
	config.LookupdPollInterval = 1 * time.Second

	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		panic(err)
	}
	consumer.AddConcurrentHandlers(nsq.HandlerFunc(handle), concurrency)
	err = consumer.ConnectToNSQD(conf.NSQIP)
	if err != nil {
		panic(err)
	}
}

// StartNsqConsume 启动nsq消费者，以后所有的消费者在这里注册
func StartNsqConsumer() {
	NsqConsumer("sync_trigger", "1", handleSyncTrigger, 20)
	NsqConsumer("message_send", "1", handleMessageSend, 20)
	NsqConsumer("message_ack", "1", handleMessageACK, 20)
	NsqConsumer("off_line", "1", handleOffLine, 20)
}

func context() *imctx.Context {
	return imctx.NewContext(db.Factoty.GetSession())
}

// handleSyncTrigger 处理消息同步出发
func handleSyncTrigger(msg *nsq.Message) error {
	var trigger transfer.SyncTrigger
	err := jsoniter.Unmarshal(msg.Body, &trigger)
	if err != nil {
		logger.Sugar.Error(err)
		return nil
	}

	ctx := context()
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

	message := transfer.Message{DeviceId: trigger.DeviceId, Type: transfer.MessageTypeSync, Messages: messages}
	produce.PublishMessage(message)

	logger.Sugar.Infow("消息同步",
		"device_id:", trigger.DeviceId,
		"user_id", trigger.UserId,
		"type", message.Type,
		"messages", message.GetLog())
	return nil
}

// handleMessageSend 处理消息发送
func handleMessageSend(msg *nsq.Message) error {
	var send transfer.MessageSend
	err := jsoniter.Unmarshal(msg.Body, &send)
	if err != nil {
		logger.Sugar.Error(err)
		return nil
	}

	send.MessageId = lib.Lid.Get()
	ctx := context()
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
		produce.PublishMessageSendACK(ack)
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
	produce.PublishMessageSendACK(ack)

	logger.Sugar.Infow("消息发送回执",
		"device_id", ack.DeviceId,
		"user_id", send.SenderUserId,
		"message_id", send.MessageId,
		"send_sequence", ack.SendSequence,
		"code", ack.Code,
	)

	return nil
}

// handleMessageACK 处理消息回执
func handleMessageACK(msg *nsq.Message) error {
	var ack transfer.MessageACK
	err := jsoniter.Unmarshal(msg.Body, &ack)
	if err != nil {
		logger.Sugar.Error(err)
		return nil
	}

	err = dao.DeviceSyncSequenceDao.UpdateSyncSequence(context(), ack.DeviceId, ack.SyncSequence)
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

// handleOffLine 处理消息离线
func handleOffLine(msg *nsq.Message) error {
	var offLine transfer.OffLine
	err := jsoniter.Unmarshal(msg.Body, &offLine)
	if err != nil {
		logger.Sugar.Error(err)
		return nil
	}

	err = dao.DeviceDao.UpdateStatus(context(), offLine.DeviceId, service.DeviceOffline)
	if err != nil {
		logger.Sugar.Error(err)
	}

	logger.Sugar.Infow("设备离线", "device_id", offLine.DeviceId, "user_id", offLine.UserId)
	return nil
}

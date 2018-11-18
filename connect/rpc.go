package connect

import (
	"goim/public/lib"
	"goim/public/logger"
	"goim/public/pb"
	"goim/public/transfer"

	"goim/public/imctx"

	"github.com/golang/protobuf/proto"
)

// LogicRPCer 逻辑层接口
type LogicRPCer interface {
	// SignIn 设备登录
	SignIn(ctx *imctx.Context, signIn transfer.SignIn) (*transfer.SignInACK, error)
	// SyncTrigger 消息同步触发
	SyncTrigger(ctx *imctx.Context, trigger transfer.SyncTrigger) error
	// MessageSend 消息发送
	MessageSend(ctx *imctx.Context, send transfer.MessageSend) error
	// MessageACK 消息投递回执
	MessageACK(ctx *imctx.Context, ack transfer.MessageACK) error
	// OffLine 下线
	OffLine(ctx *imctx.Context, deviceId int64, userId int64) error
}

var LogicRPC LogicRPCer

type connectRPC struct{}

var ConnectRPC = new(connectRPC)

// SendMessage 处理消息投递
func (*connectRPC) SendMessage(message transfer.Message) error {
	ctx := load(message.DeviceId)
	if ctx == nil {
		logger.Sugar.Error("ctx id nil")
		return nil
	}

	messages := make([]*pb.MessageItem, 0, len(message.Messages))
	for _, v := range message.Messages {
		item := new(pb.MessageItem)

		item.MessageId = v.MessageId
		item.SenderType = int32(v.SenderType)
		item.SenderId = v.SenderId
		item.SenderDeviceId = v.SenderDeviceId
		item.ReceiverType = int32(v.ReceiverType)
		item.ReceiverId = v.ReceiverId
		item.Type = int32(v.Type)
		item.Content = v.Content
		item.SyncSequence = v.Sequence
		item.SendTime = lib.UnixTime(v.SendTime)

		messages = append(messages, item)
	}

	content, err := proto.Marshal(&pb.Message{Messages: messages})
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	err = ctx.Codec.Eecode(Package{Code: CodeMessage, Content: content}, WriteDeadline)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// SendMessageSendACK 处理消息发送回执
func (*connectRPC) SendMessageSendACK(ack transfer.MessageSendACK) error {
	content, err := proto.Marshal(&pb.MessageSendACK{SendSequence: ack.SendSequence, Code: int32(ack.Code)})
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	ctx := load(ack.DeviceId)
	if ctx == nil {
		logger.Sugar.Error(err)
		return err
	}

	err = ctx.Codec.Eecode(Package{Code: CodeMessageSendACK, Content: content}, WriteDeadline)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

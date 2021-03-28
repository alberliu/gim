package service

import (
	"context"
	"gim/internal/logic/cache"
	"gim/internal/logic/model"
	"gim/pkg/gerrors"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/topic"
	"gim/pkg/util"
	"time"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

type pushService struct{}

var PushService = new(pushService)

func (s *pushService) PushToUser(ctx context.Context, userId int64, code pb.PushCode, message proto.Message, isPersist bool) error {
	logger.Logger.Debug("push_to_user",
		zap.Int64("request_id", grpclib.GetCtxRequstId(ctx)),
		zap.Int64("user_id", userId),
		zap.Int32("code", int32(code)),
		zap.Any("message", message))

	messageBuf, err := proto.Marshal(message)
	if err != nil {
		return gerrors.WrapError(err)
	}

	commandBuf, err := proto.Marshal(&pb.Command{Code: int32(code), Data: messageBuf})
	if err != nil {
		return gerrors.WrapError(err)
	}

	_, err = MessageService.SendToUser(ctx,
		model.Sender{
			SenderType: pb.SenderType_ST_SYSTEM,
			SenderId:   0,
			DeviceId:   0,
		},
		userId,
		pb.SendMessageReq{
			ReceiverType:   pb.ReceiverType_RT_USER,
			ReceiverId:     userId,
			ToUserIds:      nil,
			MessageType:    pb.MessageType_MT_COMMAND,
			MessageContent: commandBuf,
			SendTime:       util.UnixMilliTime(time.Now()),
			IsPersist:      isPersist,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *pushService) PushToGroup(ctx context.Context, groupId int64, code pb.PushCode, message proto.Message, isPersist bool) error {
	logger.Logger.Debug("push_to_group",
		zap.Int64("request_id", grpclib.GetCtxRequstId(ctx)),
		zap.Int64("group_id", groupId),
		zap.Int32("code", int32(code)),
		zap.Any("message", message))

	messageBuf, err := proto.Marshal(message)
	if err != nil {
		return gerrors.WrapError(err)
	}

	commandBuf, err := proto.Marshal(&pb.Command{Code: int32(code), Data: messageBuf})
	if err != nil {
		return gerrors.WrapError(err)
	}

	_, err = MessageService.SendToGroup(ctx,
		model.Sender{
			SenderType: pb.SenderType_ST_SYSTEM,
			SenderId:   0,
			DeviceId:   0,
		},
		pb.SendMessageReq{
			ReceiverType:   pb.ReceiverType_RT_GROUP,
			ReceiverId:     groupId,
			ToUserIds:      nil,
			MessageType:    pb.MessageType_MT_COMMAND,
			MessageContent: commandBuf,
			SendTime:       util.UnixMilliTime(time.Now()),
			IsPersist:      isPersist,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *pushService) PushRoom(ctx context.Context, sender model.Sender, req *pb.PushRoomReq) error {
	MessageService.AddInfo(&sender)

	msg := pb.PushRoomMsg{
		RoomId: req.RoomId,
		MessageSend: &pb.MessageSend{
			Message: &pb.Message{
				Sender:         model.SenderToPB(sender),
				ReceiverType:   pb.ReceiverType_RT_ROOM,
				ReceiverId:     req.RoomId,
				ToUserIds:      nil,
				MessageType:    req.MessageType,
				MessageContent: req.MessageContent,
				Seq:            0,
				SendTime:       util.UnixMilliTime(time.Now()),
				Status:         0,
			},
		},
	}
	bytes, err := proto.Marshal(&msg)
	if err != nil {
		return gerrors.WrapError(err)
	}
	cache.Queue.Publish(topic.PushAllQueue, bytes)
	return nil
}

func (s *pushService) PushAll(ctx context.Context, req *pb.PushAllReq) error {
	msg := pb.PushAllMsg{
		MessageSend: &pb.MessageSend{
			Message: &pb.Message{
				Sender:         &pb.Sender{SenderType: pb.SenderType_ST_BUSINESS},
				ReceiverType:   pb.ReceiverType_RT_ROOM,
				ToUserIds:      nil,
				MessageType:    req.MessageType,
				MessageContent: req.MessageContent,
				Seq:            0,
				SendTime:       util.UnixMilliTime(time.Now()),
				Status:         0,
			},
		},
	}
	bytes, err := proto.Marshal(&msg)
	if err != nil {
		return gerrors.WrapError(err)
	}
	cache.Queue.Publish(topic.PushAllQueue, bytes)
	return nil
}

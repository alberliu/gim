package service

import (
	"context"
	"gim/internal/logic/cache"
	"gim/internal/logic/model"
	"gim/pkg/gerrors"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc"
	"gim/pkg/topic"
	"gim/pkg/util"
	"time"

	"github.com/golang/protobuf/proto"
)

type roomService struct{}

var RoomService = new(roomService)

func (s *roomService) Push(ctx context.Context, sender model.Sender, req *pb.PushRoomReq) error {
	MessageService.AddSenderInfo(&sender)

	seq, err := SeqService.GetRoomNext(ctx, req.RoomId)
	if err != nil {
		return err
	}

	msg := &pb.Message{
		Sender:         model.SenderToPB(sender),
		ReceiverType:   pb.ReceiverType_RT_ROOM,
		ReceiverId:     req.RoomId,
		ToUserIds:      nil,
		MessageType:    req.MessageType,
		MessageContent: req.MessageContent,
		Seq:            seq,
		SendTime:       util.UnixMilliTime(time.Now()),
		Status:         0,
	}

	if req.IsPersist {
		err = s.AddMessage(req.RoomId, *msg)
		if err != nil {
			return err
		}
	}

	pushRoomMsg := pb.PushRoomMsg{
		RoomId: req.RoomId,
		MessageSend: &pb.MessageSend{
			Message: msg,
		},
	}
	bytes, err := proto.Marshal(&pushRoomMsg)
	if err != nil {
		return gerrors.WrapError(err)
	}
	cache.Queue.Publish(topic.PushRoomTopic, bytes)
	return nil
}

func (s *roomService) AddMessage(roomId int64, msg pb.Message) error {
	err := cache.RoomMessageCache.Add(roomId, msg)
	if err != nil {
		return err
	}
	return s.DelExpireMessage(roomId)
}

// DelExpireMessage 删除过期消息
func (s *roomService) DelExpireMessage(roomId int64) error {
	var (
		index int64 = 0
		stop  bool
		min   int64
		max   int64
	)

	for {
		msgs, err := cache.RoomMessageCache.ListByIndex(roomId, index, index+20)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			break
		}

		for _, v := range msgs {
			if v.SendTime > util.UnixMilliTime(time.Now().Add(-cache.RoomMessageExpireTime)) {
				stop = true
				break
			}

			if min == 0 {
				min = v.Seq
			}
			max = v.Seq
		}
		if stop {
			break
		}
	}

	return cache.RoomMessageCache.DelBySeq(roomId, min, max)
}

// DelExpireMessage 删除过期消息
func (s *roomService) SubscribeRoom(ctx context.Context, req pb.SubscribeRoomReq) error {
	msgs, err := cache.RoomMessageCache.List(req.RoomId, req.Seq)
	if err != nil {
		return err
	}

	for i := range msgs {
		_, err := rpc.ConnectIntClient.DeliverMessage(grpclib.ContextWithAddr(ctx, req.ConnAddr), &pb.DeliverMessageReq{
			DeviceId: req.DeviceId,
			MessageSend: &pb.MessageSend{
				Message: msgs[i],
			},
		})
		if err != nil {
			logger.Sugar.Error(err)
		}
	}
	return nil
}

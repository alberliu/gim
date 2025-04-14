package room

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	"gim/pkg/gerrors"
	"gim/pkg/grpclib/picker"
	"gim/pkg/logger"
	"gim/pkg/mq"
	"gim/pkg/protocol/pb"
	"gim/pkg/rpc"
	"gim/pkg/util"
)

type service struct{}

var Service = new(service)

func (s *service) Push(ctx context.Context, req *pb.PushRoomReq) error {
	seq, err := SeqRepo.GetNextSeq(req.RoomId)
	if err != nil {
		return err
	}

	msg := &pb.Message{
		Code:     req.Code,
		Content:  req.Content,
		Seq:      seq,
		SendTime: util.UnixMilliTime(time.Now()),
	}
	if req.IsPersist {
		err = s.AddMessage(req.RoomId, msg)
		if err != nil {
			return err
		}
	}

	pushRoomMsg := pb.PushRoomMsg{
		RoomId:  req.RoomId,
		Message: msg,
	}
	bytes, err := proto.Marshal(&pushRoomMsg)
	if err != nil {
		return gerrors.WrapError(err)
	}
	var topicName = mq.PushRoomTopic
	if req.IsPriority {
		topicName = mq.PushRoomPriorityTopic
	}
	err = mq.Publish(topicName, bytes)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddMessage(roomId int64, msg *pb.Message) error {
	err := MessageRepo.Add(roomId, msg)
	if err != nil {
		return err
	}
	return s.DelExpireMessage(roomId)
}

// DelExpireMessage 删除过期消息
func (s *service) DelExpireMessage(roomId int64) error {
	var (
		index int64 = 0
		stop  bool
		min   int64
		max   int64
	)

	for {
		msgs, err := MessageRepo.ListByIndex(roomId, index, index+20)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			break
		}

		for _, v := range msgs {
			if v.SendTime > util.UnixMilliTime(time.Now().Add(-MessageExpireTime)) {
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

	return MessageRepo.DelBySeq(roomId, min, max)
}

// SubscribeRoom 订阅房间
func (s *service) SubscribeRoom(ctx context.Context, req *pb.SubscribeRoomReq) error {
	if req.Seq == 0 {
		return nil
	}

	messages, err := MessageRepo.List(req.RoomId, req.Seq)
	if err != nil {
		return err
	}

	for i := range messages {
		_, err := rpc.GetConnectIntClient().DeliverMessage(picker.ContextWithAddr(ctx, req.ConnAddr), &pb.DeliverMessageReq{
			DeviceId: req.DeviceId,
			Message:  messages[i],
		})
		if err != nil {
			logger.Sugar.Error(err)
		}
	}
	return nil
}

package room

import (
	"context"
	"gim/pkg/gerrors"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/mq"
	"gim/pkg/pb"
	"gim/pkg/rpc"
	"gim/pkg/util"
	"time"

	"google.golang.org/protobuf/proto"
)

type roomService struct{}

var RoomService = new(roomService)

func (s *roomService) Push(ctx context.Context, sender *pb.Sender, req *pb.PushRoomReq) error {
	s.AddSenderInfo(sender)

	seq, err := RoomSeqRepo.GetNextSeq(req.RoomId)
	if err != nil {
		return err
	}

	msg := &pb.Message{
		Sender:         sender,
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
		err = s.AddMessage(req.RoomId, msg)
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

func (s *roomService) AddMessage(roomId int64, msg *pb.Message) error {
	err := RoomMessageRepo.Add(roomId, msg)
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
		msgs, err := RoomMessageRepo.ListByIndex(roomId, index, index+20)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			break
		}

		for _, v := range msgs {
			if v.SendTime > util.UnixMilliTime(time.Now().Add(-RoomMessageExpireTime)) {
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

	return RoomMessageRepo.DelBySeq(roomId, min, max)
}

// SubscribeRoom 订阅房间
func (s *roomService) SubscribeRoom(ctx context.Context, req *pb.SubscribeRoomReq) error {
	if req.Seq == 0 {
		return nil
	}

	messages, err := RoomMessageRepo.List(req.RoomId, req.Seq)
	if err != nil {
		return err
	}

	for i := range messages {
		_, err := rpc.GetConnectIntClient().DeliverMessage(grpclib.ContextWithAddr(ctx, req.ConnAddr), &pb.DeliverMessageReq{
			DeviceId: req.DeviceId,
			MessageSend: &pb.MessageSend{
				Message: messages[i],
			},
		})
		if err != nil {
			logger.Sugar.Error(err)
		}
	}
	return nil
}

func (*roomService) AddSenderInfo(sender *pb.Sender) {
	if sender.SenderType == pb.SenderType_ST_USER {
		user, err := rpc.GetBusinessIntClient().GetUser(context.TODO(), &pb.GetUserReq{UserId: sender.SenderId})
		if err == nil && user != nil {
			sender.AvatarUrl = user.User.AvatarUrl
			sender.Nickname = user.User.Nickname
			sender.Extra = user.User.Extra
		}
	}
}

package room

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	"gim/pkg/mq"
	"gim/pkg/protocol/pb/connectpb"
	pb "gim/pkg/protocol/pb/logicpb"
	"gim/pkg/util"
)

type app struct{}

var App = new(app)

// Push 推送房间消息
func (s *app) Push(ctx context.Context, req *pb.PushRoomRequest) error {
	seq, err := SeqRepo.GetNextSeq(req.RoomId)
	if err != nil {
		return err
	}

	msg := &pb.Message{
		Code:      req.Code,
		Content:   req.Content,
		Seq:       seq,
		CreatedAt: util.UnixMilliTime(time.Now()),
	}
	if req.IsPersist {
		err = s.addMessage(req.RoomId, msg)
		if err != nil {
			return err
		}
	}

	pushRoomMsg := connectpb.PushRoomMsg{
		RoomId:  req.RoomId,
		Message: msg,
	}
	buf, err := proto.Marshal(&pushRoomMsg)
	if err != nil {
		return err
	}
	var topicName = mq.PushRoomTopic
	if req.IsPriority {
		topicName = mq.PushRoomPriorityTopic
	}
	return mq.Publish(topicName, buf)
}

// SubscribeRoom 订阅房间
func (s *app) SubscribeRoom(ctx context.Context, req *pb.SubscribeRoomRequest) error {
	return nil
}

func (s *app) addMessage(roomId uint64, msg *pb.Message) error {
	err := MessageRepo.Add(roomId, msg)
	if err != nil {
		return err
	}
	return s.delExpireMessage(roomId)
}

// DelExpireMessage 删除过期消息
func (s *app) delExpireMessage(roomId uint64) error {
	var (
		index int64 = 0
		stop  bool
		min   uint64
		max   uint64
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
			if v.CreatedAt > util.UnixMilliTime(time.Now().Add(-MessageExpireTime)) {
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

package room

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	"gim/pkg/mq"
	"gim/pkg/protocol/pb/connectpb"
	pb "gim/pkg/protocol/pb/logicpb"
)

type app struct{}

var App = new(app)

// Push 推送房间消息
func (s *app) Push(ctx context.Context, req *pb.PushRoomRequest) error {
	msg := &connectpb.Message{
		Command:   req.Command,
		Content:   req.Content,
		CreatedAt: time.Now().Unix(),
		RoomId:    req.RoomId,
	}

	pushRoomMsg := connectpb.PushRoomMessage{
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

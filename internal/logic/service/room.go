package service

import (
	"context"
	"gim/internal/logic/cache"
	"gim/internal/logic/model"
	"gim/pkg/gerrors"
	"gim/pkg/pb"
	"gim/pkg/topic"
	"gim/pkg/util"
	"time"

	"github.com/golang/protobuf/proto"
)

type roomService struct{}

var RoomService = new(roomService)

func (s *roomService) Push(ctx context.Context, sender model.Sender, req *pb.PushRoomReq) error {
	MessageService.AddSenderInfo(&sender)

	// 需要将消息保存

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
	cache.Queue.Publish(topic.PushRoomTopic, bytes)
	return nil
}

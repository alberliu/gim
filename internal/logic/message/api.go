package message

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"gim/pkg/protocol/pb/connectpb"
	pb "gim/pkg/protocol/pb/logicpb"
)

type MessageIntService struct {
	pb.UnsafeMessageIntServiceServer
}

func (m MessageIntService) Sync(ctx context.Context, request *pb.SyncRequest) (*connectpb.SyncReply, error) {
	return App.Sync(ctx, request.UserId, request.Seq)
}

// MessageACK 设备收到消息ack
func (*MessageIntService) MessageACK(ctx context.Context, request *pb.MessageACKRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, App.MessageAck(ctx, request.UserId, request.DeviceId, request.DeviceAck)
}

// PushToUsers 推送
func (*MessageIntService) PushToUsers(ctx context.Context, request *pb.PushToUsersRequest) (*pb.PushToUsersReply, error) {
	messageID, err := App.PushToUsersWithCommand(ctx, request.UserIds, request.Command, request.Content, request.IsPersist)
	if err != nil {
		return nil, err
	}
	return &pb.PushToUsersReply{MessageId: messageID}, nil
}

// PushToAll 全服推送
func (s *MessageIntService) PushToAll(ctx context.Context, request *pb.PushToAllRequest) (*emptypb.Empty, error) {
	err := App.PushToAll(ctx, request)
	return &emptypb.Empty{}, err
}

package message

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "gim/pkg/protocol/pb/logicpb"
)

type MessageIntService struct {
	pb.UnsafeMessageIntServiceServer
}

// Sync 设备同步消息
func (*MessageIntService) Sync(ctx context.Context, request *pb.SyncRequest) (*pb.SyncReply, error) {
	return App.Sync(ctx, request.UserId, request.Seq)
}

// MessageACK 设备收到消息ack
func (*MessageIntService) MessageACK(ctx context.Context, request *pb.MessageACKRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, App.MessageAck(ctx, request.UserId, request.DeviceId, request.DeviceAck)
}

// Pushs 推送
func (*MessageIntService) Pushs(ctx context.Context, request *pb.PushsRequest) (*pb.PushsReply, error) {
	messageID, err := App.PushToUserData(ctx, request.UserIds, request.Code, request.Content, request.IsPersist)
	if err != nil {
		return nil, err
	}
	return &pb.PushsReply{MessageId: messageID}, nil
}

// PushAll 全服推送
func (s *MessageIntService) PushAll(ctx context.Context, request *pb.PushAllRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, App.PushAll(ctx, request)
}

package message

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"gim/pkg/md"
	pb "gim/pkg/protocol/pb/logicpb"
)

type MessageExtService struct {
	pb.UnsafeMessageExtServiceServer
}

func (m MessageExtService) Sync(ctx context.Context, request *pb.SyncRequest) (*pb.SyncReply, error) {
	userID, _, _ := md.GetData(ctx)
	return App.Sync(ctx, userID, request.Seq)
}

type MessageIntService struct {
	pb.UnsafeMessageIntServiceServer
}

// MessageACK 设备收到消息ack
func (*MessageIntService) MessageACK(ctx context.Context, request *pb.MessageACKRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, App.MessageAck(ctx, request.UserId, request.DeviceId, request.DeviceAck)
}

// Pushs 推送
func (*MessageIntService) Pushs(ctx context.Context, request *pb.PushsRequest) (*pb.PushsReply, error) {
	messageID, err := App.PushContent(ctx, request.UserIds, request.Command, request.Content, request.IsPersist)
	if err != nil {
		return nil, err
	}
	return &pb.PushsReply{MessageId: messageID}, nil
}

// PushAll 全服推送
func (s *MessageIntService) PushAll(ctx context.Context, request *pb.PushAllRequest) (*emptypb.Empty, error) {
	err := App.PushAll(ctx, request)
	return &emptypb.Empty{}, err
}

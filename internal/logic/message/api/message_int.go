package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"gim/internal/logic/message/app"
	"gim/pkg/protocol/pb/connectpb"
	pb "gim/pkg/protocol/pb/logicpb"
)

type MessageIntService struct {
	pb.UnsafeMessageIntServiceServer
}

// PushToUsers 推送
func (*MessageIntService) PushToUsers(ctx context.Context, request *pb.PushToUsersRequest) (*pb.PushToUsersReply, error) {
	message := &connectpb.Message{
		Command: request.Command,
		Content: request.Content,
	}
	messageID, err := app.MessageApp.PushToUsers(ctx, request.UserIds, message, request.IsPersist)
	if err != nil {
		return nil, err
	}
	return &pb.PushToUsersReply{MessageId: messageID}, nil
}

// PushToAll 全服推送
func (s *MessageIntService) PushToAll(ctx context.Context, request *pb.PushToAllRequest) (*emptypb.Empty, error) {
	err := app.MessageApp.PushToAll(ctx, request)
	return &emptypb.Empty{}, err
}

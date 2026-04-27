package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"gim/internal/logic/message/app"
	pb "gim/pkg/protocol/pb/logicpb"
)

type MessageIntService struct {
	pb.UnsafeMessageIntServiceServer
}

// PushToUsers 推送
func (*MessageIntService) PushToUsers(ctx context.Context, request *pb.PushToUsersRequest) (*pb.PushToUsersReply, error) {
	return app.MessageApp.PushToUsers(ctx, request)
}

// PushToAll 全服推送
func (s *MessageIntService) PushToAll(ctx context.Context, request *pb.PushToAllRequest) (*emptypb.Empty, error) {
	err := app.MessageApp.PushToAll(ctx, request)
	return &emptypb.Empty{}, err
}

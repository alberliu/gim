package room

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "gim/pkg/protocol/pb/logicpb"
)

type RoomIntService struct {
	pb.UnsafeRoomIntServiceServer
}

// SubscribeRoom 订阅房间
func (s *RoomIntService) SubscribeRoom(ctx context.Context, request *pb.SubscribeRoomRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, App.SubscribeRoom(ctx, request)
}

// PushRoom 推送房间
func (s *RoomIntService) PushRoom(ctx context.Context, request *pb.PushRoomRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, App.Push(ctx, request)
}

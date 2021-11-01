package app

import (
	"context"
	"gim/internal/logic/domain/room"
	"gim/pkg/pb"
)

type roomApp struct{}

var RoomApp = new(roomApp)

// Push 推送房间消息
func (s *roomApp) Push(ctx context.Context, sender *pb.Sender, req *pb.PushRoomReq) error {
	return room.RoomService.Push(ctx, sender, req)
}

// SubscribeRoom 订阅房间
func (s *roomApp) SubscribeRoom(ctx context.Context, req *pb.SubscribeRoomReq) error {
	return room.RoomService.SubscribeRoom(ctx, req)
}

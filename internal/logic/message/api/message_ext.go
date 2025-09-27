package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"gim/internal/logic/message/app"
	"gim/pkg/md"
	pb "gim/pkg/protocol/pb/logicpb"
)

type MessageExtService struct {
	pb.UnsafeMessageExtServiceServer
}

func (*MessageExtService) Sync(ctx context.Context, request *pb.SyncRequest) (*pb.SyncReply, error) {
	userID := md.GetUserID(ctx)
	return app.MessageApp.Sync(ctx, userID, request.Seq)
}

// MessageACK 设备收到消息ack
func (*MessageExtService) MessageACK(ctx context.Context, request *pb.MessageACKRequest) (*emptypb.Empty, error) {
	userID := md.GetUserID(ctx)
	deviceID := md.GetDeviceID(ctx)
	return &emptypb.Empty{}, app.DeviceACKApp.MessageAck(ctx, userID, deviceID, request.DeviceAck)
}

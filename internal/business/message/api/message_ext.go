package api

import (
	"context"

	"gim/internal/business/message/app"
	pb "gim/pkg/gen/proto/businesspb"
	"gim/pkg/md"
)

type MessageExtService struct {
	pb.UnsafeMessageExtServiceServer
}

func (*MessageExtService) SendFriendMessage(ctx context.Context, req *pb.SendFriendMessageRequest) (*pb.SendFriendMessageReply, error) {
	userID := md.GetUserID(ctx)
	deviceID := md.GetDeviceID(ctx)

	reply, err := app.MessageApp.SendToFriend(ctx, deviceID, userID, req)
	if err != nil {
		return nil, err
	}
	return &pb.SendFriendMessageReply{
		MessageId:   reply.MessageId,
		FromUserSeq: reply.FromUserSeq,
	}, nil
}

func (*MessageExtService) SendGroupMessage(ctx context.Context, req *pb.SendGroupMessageRequest) (*pb.SendGroupMessageReply, error) {
	userID := md.GetUserID(ctx)
	deviceID := md.GetDeviceID(ctx)

	reply, err := app.MessageApp.SendToGroup(ctx, deviceID, userID, req)
	if err != nil {
		return nil, err
	}
	return &pb.SendGroupMessageReply{
		MessageId:   reply.MessageId,
		FromUserSeq: reply.FromUserSeq,
	}, nil
}

func (*MessageExtService) SendRoomMessage(ctx context.Context, req *pb.SendRoomMessageRequest) (*pb.SendRoomMessageReply, error) {
	err := app.MessageApp.SendToRoom(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.SendRoomMessageReply{}, nil
}

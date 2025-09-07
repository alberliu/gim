package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"gim/internal/business/friend/app"
	"gim/pkg/md"
	pb "gim/pkg/protocol/pb/businesspb"
)

type FriendExtService struct {
	pb.UnsafeFriendExtServiceServer
}

// SendMessage 发送好友消息
func (*FriendExtService) SendMessage(ctx context.Context, request *pb.SendFriendMessageRequest) (*pb.SendFriendMessageReply, error) {
	userID := md.GetUserID(ctx)
	deviceID := md.GetDeviceID(ctx)

	messageId, err := app.FriendApp.SendToFriend(ctx, deviceID, userID, request)
	if err != nil {
		return nil, err
	}
	return &pb.SendFriendMessageReply{MessageId: messageId}, nil
}

func (s *FriendExtService) Add(ctx context.Context, request *pb.FriendAddRequest) (*emptypb.Empty, error) {
	userID := md.GetUserID(ctx)

	err := app.FriendApp.AddFriend(ctx, userID, request.FriendId, request.Remarks, request.Description)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *FriendExtService) Agree(ctx context.Context, request *pb.FriendAgreeRequest) (*emptypb.Empty, error) {
	userID := md.GetUserID(ctx)

	err := app.FriendApp.AgreeAddFriend(ctx, userID, request.UserId, request.Remarks)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *FriendExtService) Set(ctx context.Context, request *pb.FriendSetRequest) (*pb.FriendSetReply, error) {
	userID := md.GetUserID(ctx)

	err := app.FriendApp.SetFriend(ctx, userID, request)
	if err != nil {
		return nil, err
	}
	return &pb.FriendSetReply{}, nil
}

func (s *FriendExtService) GetFriends(ctx context.Context, request *emptypb.Empty) (*pb.GetFriendsReply, error) {
	userId := md.GetUserID(ctx)
	friends, err := app.FriendApp.List(ctx, userId)
	return &pb.GetFriendsReply{Friends: friends}, err
}

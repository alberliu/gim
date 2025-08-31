package friend

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"gim/pkg/md"
	pb "gim/pkg/protocol/pb/businesspb"
)

type FriendExtService struct {
	pb.UnsafeFriendExtServiceServer
}

// SendMessage 发送好友消息
func (*FriendExtService) SendMessage(ctx context.Context, request *pb.SendFriendMessageRequest) (*pb.SendFriendMessageReply, error) {
	userID, deviceID, err := md.GetData(ctx)
	if err != nil {
		return nil, err
	}

	messageId, err := App.SendToFriend(ctx, deviceID, userID, request)
	if err != nil {
		return nil, err
	}
	return &pb.SendFriendMessageReply{MessageId: messageId}, nil
}

func (s *FriendExtService) Add(ctx context.Context, request *pb.FriendAddRequest) (*emptypb.Empty, error) {
	userID, _, err := md.GetData(ctx)
	if err != nil {
		return nil, err
	}

	err = App.AddFriend(ctx, userID, request.FriendId, request.Remarks, request.Description)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *FriendExtService) Agree(ctx context.Context, request *pb.FriendAgreeRequest) (*emptypb.Empty, error) {
	userID, _, err := md.GetData(ctx)
	if err != nil {
		return nil, err
	}

	err = App.AgreeAddFriend(ctx, userID, request.UserId, request.Remarks)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *FriendExtService) Set(ctx context.Context, request *pb.FriendSetRequest) (*pb.FriendSetReply, error) {
	userID, _, err := md.GetData(ctx)
	if err != nil {
		return nil, err
	}

	err = App.SetFriend(ctx, userID, request)
	if err != nil {
		return nil, err
	}
	return &pb.FriendSetReply{}, nil
}

func (s *FriendExtService) GetFriends(ctx context.Context, request *emptypb.Empty) (*pb.GetFriendsReply, error) {
	userId, _, err := md.GetData(ctx)
	if err != nil {
		return nil, err
	}
	friends, err := App.List(ctx, userId)
	return &pb.GetFriendsReply{Friends: friends}, err
}

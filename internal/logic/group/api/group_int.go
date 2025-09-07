package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"gim/internal/logic/group/app"
	pb "gim/pkg/protocol/pb/logicpb"
)

type GroupIntService struct {
	pb.UnsafeGroupIntServiceServer
}

// Push 发送群组消息
func (*GroupIntService) Push(ctx context.Context, request *pb.GroupPushRequest) (*pb.GroupPushReply, error) {
	messageID, err := app.GroupApp.Push(ctx, request)
	if err != nil {
		return nil, err
	}
	return &pb.GroupPushReply{MessageId: messageID}, nil
}

// Create 创建群组
func (*GroupIntService) Create(ctx context.Context, request *pb.GroupCreateRequest) (*pb.GroupCreateReply, error) {
	groupID, err := app.GroupApp.Create(ctx, request.Group)
	return &pb.GroupCreateReply{GroupId: groupID}, err
}

// Update 更新群组
func (*GroupIntService) Update(ctx context.Context, request *pb.GroupUpdateRequest) (*emptypb.Empty, error) {
	err := app.GroupApp.Update(ctx, request.Group)
	return &emptypb.Empty{}, err
}

// Get 获取群组信息
func (*GroupIntService) Get(ctx context.Context, request *pb.GroupGetRequest) (*pb.GroupGetReply, error) {
	group, err := app.GroupApp.Get(ctx, request.GroupId)
	return &pb.GroupGetReply{Group: group}, err
}

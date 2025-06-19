package group

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"gim/pkg/md"
	pb "gim/pkg/protocol/pb/logicpb"
)

type GroupExtService struct {
	pb.UnsafeGroupExtServiceServer
}

// SendMessage 发送群组消息
func (*GroupExtService) SendMessage(ctx context.Context, request *pb.SendGroupMessageRequest) (*pb.SendGroupMessageReply, error) {
	userID, deviceID, err := md.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	messageID, err := App.SendMessage(ctx, deviceID, userID, request)
	if err != nil {
		return nil, err
	}
	return &pb.SendGroupMessageReply{MessageId: messageID}, nil
}

// Create 创建群组
func (*GroupExtService) Create(ctx context.Context, request *pb.GroupCreateRequest) (*pb.GroupCreateReply, error) {
	userID, _, err := md.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	groupID, err := App.CreateGroup(ctx, userID, request)
	return &pb.GroupCreateReply{GroupId: groupID}, err
}

// Update 更新群组
func (*GroupExtService) Update(ctx context.Context, request *pb.GroupUpdateRequest) (*emptypb.Empty, error) {
	userID, _, err := md.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	err = App.Update(ctx, userID, request)
	return &emptypb.Empty{}, err
}

// Get 获取群组信息
func (*GroupExtService) Get(ctx context.Context, request *pb.GroupGetRequest) (*pb.GroupGetReply, error) {
	group, err := App.GetGroup(ctx, request.GroupId)
	return &pb.GroupGetReply{Group: group}, err
}

// List 获取用户加入的所有群组
func (*GroupExtService) List(ctx context.Context, in *emptypb.Empty) (*pb.GroupListReply, error) {
	userID, _, err := md.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := App.GetUserGroups(ctx, userID)
	return &pb.GroupListReply{Groups: groups}, err
}

func (s *GroupExtService) AddMembers(ctx context.Context, in *pb.AddMembersRequest) (*pb.AddMembersReply, error) {
	userID, _, err := md.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	err = App.AddMembers(ctx, userID, in.GroupId, in.UserIds)
	return &pb.AddMembersReply{}, err
}

// UpdateMember 更新群组成员信息
func (*GroupExtService) UpdateMember(ctx context.Context, in *pb.UpdateMemberRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, App.UpdateMember(ctx, in)
}

// DeleteMember 添加群组成员
func (*GroupExtService) DeleteMember(ctx context.Context, in *pb.DeleteMemberRequest) (*emptypb.Empty, error) {
	userID, _, err := md.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	err = App.DeleteMember(ctx, in.GroupId, in.UserId, userID)
	return &emptypb.Empty{}, err
}

// GetMembers 获取群组成员信息
func (s *GroupExtService) GetMembers(ctx context.Context, in *pb.GetMembersRequest) (*pb.GetMembersReply, error) {
	members, err := App.GetMembers(ctx, in.GroupId)
	return &pb.GetMembersReply{Members: members}, err
}

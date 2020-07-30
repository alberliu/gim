package api

import (
	"context"
	"gim/internal/logic/model"
	"gim/internal/logic/service"
	"gim/pkg/gerrors"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/pb"
)

type LogicExtServer struct{}

// RegisterDevice 注册设备
func (*LogicExtServer) RegisterDevice(ctx context.Context, in *pb.RegisterDeviceReq) (*pb.RegisterDeviceResp, error) {
	device := model.Device{
		Type:          in.Type,
		Brand:         in.Brand,
		Model:         in.Model,
		SystemVersion: in.SystemVersion,
		SDKVersion:    in.SdkVersion,
	}

	if device.Type == 0 || device.Brand == "" || device.Model == "" ||
		device.SystemVersion == "" || device.SDKVersion == "" {
		return nil, gerrors.ErrBadRequest
	}

	id, err := service.DeviceService.Register(ctx, device)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterDeviceResp{DeviceId: id}, nil
}

// SendMessage 发送消息
func (*LogicExtServer) SendMessage(ctx context.Context, in *pb.SendMessageReq) (*pb.SendMessageResp, error) {
	userId, deviceId, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	sender := model.Sender{
		SenderType: pb.SenderType_ST_USER,
		SenderId:   userId,
		DeviceId:   deviceId,
	}
	seq, err := service.MessageService.Send(ctx, sender, *in)
	if err != nil {
		return nil, err
	}
	return &pb.SendMessageResp{Seq: seq}, nil
}

func (s *LogicExtServer) AddFriend(ctx context.Context, in *pb.AddFriendReq) (*pb.AddFriendResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	err = service.FriendService.AddFriend(ctx, userId, in.FriendId, in.Remarks, in.Description)
	if err != nil {
		return nil, err
	}

	return &pb.AddFriendResp{}, nil
}

func (s *LogicExtServer) AgreeAddFriend(ctx context.Context, in *pb.AgreeAddFriendReq) (*pb.AgreeAddFriendResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	err = service.FriendService.AgreeAddFriend(ctx, userId, in.UserId, in.Remarks)
	if err != nil {
		return nil, err
	}

	return &pb.AgreeAddFriendResp{}, nil
}

func (s *LogicExtServer) GetFriends(ctx context.Context, in *pb.GetFriendsReq) (*pb.GetFriendsResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	friends, err := service.FriendService.List(ctx, userId)
	return &pb.GetFriendsResp{Friends: friends}, err
}

// CreateGroup 创建群组
func (*LogicExtServer) CreateGroup(ctx context.Context, in *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {
	groupId, err := service.GroupService.Create(ctx, model.Group{
		Name:         in.Name,
		Introduction: in.Introduction,
		Type:         in.Type,
		Extra:        in.Extra,
	})
	return &pb.CreateGroupResp{GroupId: groupId}, err
}

// UpdateGroup 更新群组
func (*LogicExtServer) UpdateGroup(ctx context.Context, in *pb.UpdateGroupReq) (*pb.UpdateGroupResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateGroupResp{}, service.GroupService.Update(ctx, userId, model.Group{
		Id:           in.GroupId,
		Name:         in.Name,
		Introduction: in.Introduction,
		Extra:        in.Extra,
	})
}

// GetGroup 获取群组信息
func (*LogicExtServer) GetGroup(ctx context.Context, in *pb.GetGroupReq) (*pb.GetGroupResp, error) {
	group, err := service.GroupService.Get(ctx, in.GroupId)
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, gerrors.ErrGroupNotExist
	}

	return &pb.GetGroupResp{
		Group: &pb.Group{
			GroupId:      group.Id,
			Name:         group.Name,
			Introduction: group.Introduction,
			UserMum:      group.UserNum,
			Type:         group.Type,
			Extra:        group.Extra,
			CreateTime:   group.CreateTime.Unix(),
			UpdateTime:   group.UpdateTime.Unix(),
		},
	}, nil
}

// GetUserGroups 获取用户加入的所有群组
func (*LogicExtServer) GetUserGroups(ctx context.Context, in *pb.GetUserGroupsReq) (*pb.GetUserGroupsResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := service.SmallGroupUserService.ListByUserId(ctx, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	pbGroups := make([]*pb.Group, 0, len(groups))
	for i := range groups {
		pbGroups = append(pbGroups, &pb.Group{
			GroupId:      groups[i].Id,
			Name:         groups[i].Name,
			Introduction: groups[i].Introduction,
			UserMum:      groups[i].UserNum,
			Type:         groups[i].Type,
			Extra:        groups[i].Extra,
			CreateTime:   groups[i].CreateTime.Unix(),
			UpdateTime:   groups[i].UpdateTime.Unix(),
		})
	}
	return &pb.GetUserGroupsResp{Groups: pbGroups}, err
}

func (s *LogicExtServer) AddGroupMembers(ctx context.Context, in *pb.AddGroupMembersReq) (*pb.AddGroupMembersResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	userIds, err := service.GroupUserService.AddUsers(ctx, userId, in.GroupId, in.UserIds)
	return &pb.AddGroupMembersResp{UserIds: userIds}, err
}

// UpdateGroupMember 更新群组成员信息
func (*LogicExtServer) UpdateGroupMember(ctx context.Context, in *pb.UpdateGroupMemberReq) (*pb.UpdateGroupMemberResp, error) {
	err := service.GroupUserService.UpdateUser(ctx, in.GroupId, in.UserId, in.Remarks, in.Extra)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateGroupMemberResp{}, nil
}

// DeleteGroupMember 添加群组成员
func (*LogicExtServer) DeleteGroupMember(ctx context.Context, in *pb.DeleteGroupMemberReq) (*pb.DeleteGroupMemberResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	err = service.GroupUserService.DeleteUser(ctx, userId, in.GroupId, in.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteGroupMemberResp{}, nil
}

// GetGroupMembers 获取群组成员信息
func (s *LogicExtServer) GetGroupMembers(ctx context.Context, in *pb.GetGroupMembersReq) (*pb.GetGroupMembersResp, error) {
	members, err := service.GroupService.GetUsers(ctx, in.GroupId)
	return &pb.GetGroupMembersResp{Members: members}, err
}

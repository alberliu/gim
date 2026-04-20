package app

import (
	"context"

	"gim/internal/logic/group/domain"
	"gim/internal/logic/group/repo"
	messageapp "gim/internal/logic/message/app"
	"gim/pkg/protocol/pb/connectpb"
	pb "gim/pkg/protocol/pb/logicpb"
)

var GroupApp = new(groupApp)

type groupApp struct{}

// Create 创建
func (*groupApp) Create(ctx context.Context, request *pb.GroupCreateRequest) (uint64, error) {
	group := &domain.Group{
		ID:           request.Group.Id,
		Name:         request.Group.Name,
		AvatarUrl:    request.Group.AvatarUrl,
		Introduction: request.Group.Introduction,
		Extra:        request.Group.Extra,
	}
	err := repo.GroupRepo.Create(ctx, group)
	if err != nil {
		return 0, err
	}

	members := make([]domain.GroupMember, 0, len(request.Members))
	for _, member := range request.Members {
		members = append(members, domain.GroupMember{
			GroupID:  group.ID,
			UserID:   member.UserId,
			Nickname: member.Nickname,
			Type:     member.Type,
			Status:   member.Status,
			Extra:    member.Extra,
		})
	}
	err = repo.GroupMemberRepo.BatchCreate(ctx, group.ID, members)
	if err != nil {
		return 0, err
	}
	return group.ID, nil
}

// Get 获取群组信息
func (*groupApp) Get(ctx context.Context, groupID uint64) (*pb.Group, error) {
	group, err := repo.GroupRepo.Get(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return group.ToProto(), nil
}

// Update 更新群组
func (*groupApp) Update(ctx context.Context, pbgroup *pb.Group) error {
	group := &domain.Group{
		ID:           pbgroup.Id,
		Name:         pbgroup.Name,
		AvatarUrl:    pbgroup.AvatarUrl,
		Introduction: pbgroup.Introduction,
		Extra:        pbgroup.Extra,
	}
	return repo.GroupRepo.Save(ctx, group)
}

// AddMember 添加成员
func (*groupApp) AddMember(ctx context.Context, request *pb.GroupAddMemberRequest) error {
	members := make([]domain.GroupMember, 0, len(request.Members))
	for _, member := range request.Members {
		members = append(members, domain.GroupMember{
			GroupID:  request.GroupId,
			UserID:   member.UserId,
			Nickname: member.Nickname,
			Type:     member.Type,
			Status:   member.Status,
			Extra:    member.Extra,
		})
	}
	return repo.GroupMemberRepo.BatchCreate(ctx, request.GroupId, members)
}

// RemoveMember 移除成员
func (*groupApp) RemoveMember(ctx context.Context, request *pb.GroupRemoveMemberRequest) error {
	return repo.GroupMemberRepo.BatchDelete(ctx, request.GroupId, request.UserIds)
}

// Push 发送群组消息
func (*groupApp) Push(ctx context.Context, request *pb.GroupPushRequest) (uint64, error) {
	_, err := repo.GroupRepo.Get(ctx, request.GroupId)
	if err != nil {
		return 0, err
	}
	members, err := repo.GroupMemberRepo.ListByGroupID(ctx, request.GroupId)
	if err != nil {
		return 0, err
	}

	memberIDs := make([]uint64, 0, len(members))
	for _, member := range members {
		memberIDs = append(memberIDs, member.UserID)
	}

	message := &connectpb.Message{
		Command: request.Command,
		Content: request.Content,
	}
	return messageapp.MessageApp.PushToUsers(ctx, memberIDs, message, request.IsPersist)
}

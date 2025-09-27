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
func (*groupApp) Create(ctx context.Context, pbgroup *pb.Group) (uint64, error) {
	group := &domain.Group{
		ID:           pbgroup.Id,
		Name:         pbgroup.Name,
		AvatarUrl:    pbgroup.AvatarUrl,
		Introduction: pbgroup.Introduction,
		Extra:        pbgroup.Extra,
		Members:      pbgroup.Members,
	}
	err := repo.GroupRepo.Create(ctx, group)
	if err != nil {
		return 0, err
	}
	return group.ID, err
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
		Members:      pbgroup.Members,
	}
	return repo.GroupRepo.Save(ctx, group)
}

// Push 发送群组消息
func (*groupApp) Push(ctx context.Context, request *pb.GroupPushRequest) (uint64, error) {
	group, err := repo.GroupRepo.Get(ctx, request.GroupId)
	if err != nil {
		return 0, err
	}

	message := &connectpb.Message{
		Command: request.Command,
		Content: request.Content,
	}
	return messageapp.MessageApp.PushToUsers(ctx, group.Members, message, request.IsPersist)
}

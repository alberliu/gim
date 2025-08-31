package group

import (
	"context"

	"gim/internal/logic/message"
	pb "gim/pkg/protocol/pb/logicpb"
)

type app struct{}

var App = new(app)

// Create 创建
func (*app) Create(ctx context.Context, pbgroup *pb.Group) (uint64, error) {
	group := &Group{
		ID:           pbgroup.Id,
		Name:         pbgroup.Name,
		AvatarUrl:    pbgroup.AvatarUrl,
		Introduction: pbgroup.Introduction,
		Extra:        pbgroup.Extra,
		Members:      pbgroup.Members,
	}
	err := Repo.Create(group)
	if err != nil {
		return 0, err
	}
	return group.ID, err
}

// Get 获取群组信息
func (*app) Get(ctx context.Context, groupID uint64) (*pb.Group, error) {
	group, err := Repo.Get(groupID)
	if err != nil {
		return nil, err
	}
	return group.ToProto(), nil
}

// Update 更新群组
func (*app) Update(ctx context.Context, pbgroup *pb.Group) error {
	group := &Group{
		ID:           pbgroup.Id,
		Name:         pbgroup.Name,
		AvatarUrl:    pbgroup.AvatarUrl,
		Introduction: pbgroup.Introduction,
		Extra:        pbgroup.Extra,
		Members:      pbgroup.Members,
	}
	return Repo.Save(group)
}

// Push 发送群组消息
func (*app) Push(ctx context.Context, request *pb.GroupPushRequest) (uint64, error) {
	group, err := Repo.Get(request.GroupId)
	if err != nil {
		return 0, err
	}

	return message.App.PushToUsersWithCommand(ctx, group.Members, request.Command, request.Content, request.IsPersist)
}

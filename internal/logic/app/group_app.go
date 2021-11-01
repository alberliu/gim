package app

import (
	"context"
	"gim/internal/logic/domain/group/model"
	"gim/internal/logic/domain/group/repo"
	"gim/pkg/pb"
)

type groupApp struct{}

var GroupApp = new(groupApp)

// CreateGroup 创建群组
func (*groupApp) CreateGroup(ctx context.Context, userId int64, in *pb.CreateGroupReq) (int64, error) {
	group := model.CreateGroup(userId, in)
	err := repo.GroupRepo.Save(group)
	if err != nil {
		return 0, err
	}
	return group.Id, nil
}

// GetGroup 获取群组信息
func (*groupApp) GetGroup(ctx context.Context, groupId int64) (*pb.Group, error) {
	group, err := repo.GroupRepo.Get(groupId)
	if err != nil {
		return nil, err
	}

	return group.ToProto(), nil
}

// GetUserGroups 获取用户加入的群组列表
func (*groupApp) GetUserGroups(ctx context.Context, userId int64) ([]*pb.Group, error) {
	groups, err := repo.GroupUserRepo.ListByUserId(userId)
	if err != nil {
		return nil, err
	}

	pbGroups := make([]*pb.Group, len(groups))
	for i := range groups {
		pbGroups[i] = groups[i].ToProto()
	}
	return pbGroups, nil
}

// Update 更新群组
func (*groupApp) Update(ctx context.Context, userId int64, update *pb.UpdateGroupReq) error {
	group, err := repo.GroupRepo.Get(update.GroupId)
	if err != nil {
		return err
	}

	err = group.Update(ctx, userId, update)
	if err != nil {
		return err
	}

	err = repo.GroupRepo.Save(group)
	if err != nil {
		return err
	}

	err = group.PushUpdate(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}

// AddMembers 添加群组成员
func (*groupApp) AddMembers(ctx context.Context, userId, groupId int64, userIds []int64) ([]int64, error) {
	group, err := repo.GroupRepo.Get(groupId)
	if err != nil {
		return nil, err
	}
	existIds, addedIds, err := group.AddMembers(ctx, userId, userIds)
	if err != nil {
		return nil, err
	}
	err = repo.GroupRepo.Save(group)
	if err != nil {
		return nil, err
	}

	err = group.PushAddMember(ctx, userId, addedIds)
	if err != nil {
		return nil, err
	}
	return existIds, nil
}

// UpdateMember 更新群组用户
func (*groupApp) UpdateMember(ctx context.Context, in *pb.UpdateGroupMemberReq) error {
	group, err := repo.GroupRepo.Get(in.GroupId)
	if err != nil {
		return err
	}
	err = group.UpdateMember(ctx, in)
	if err != nil {
		return err
	}
	err = repo.GroupRepo.Save(group)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMember 删除群组成员
func (*groupApp) DeleteMember(ctx context.Context, groupId int64, userId int64, optId int64) error {
	group, err := repo.GroupRepo.Get(groupId)
	if err != nil {
		return err
	}
	err = group.DeleteMember(ctx, optId, userId)
	if err != nil {
		return err
	}
	err = repo.GroupRepo.Save(group)
	if err != nil {
		return err
	}

	err = group.PushDeleteMember(ctx, optId, userId)
	if err != nil {
		return err
	}
	return nil
}

// GetMembers 获取群组成员
func (*groupApp) GetMembers(ctx context.Context, groupId int64) ([]*pb.GroupMember, error) {
	group, err := repo.GroupRepo.Get(groupId)
	if err != nil {
		return nil, err
	}
	return group.GetMembers(ctx)
}

// SendMessage 获取群组成员
func (*groupApp) SendMessage(ctx context.Context, sender *pb.Sender, req *pb.SendMessageReq) (int64, error) {
	group, err := repo.GroupRepo.Get(req.ReceiverId)
	if err != nil {
		return 0, err
	}

	return group.SendMessage(ctx, sender, req)
}

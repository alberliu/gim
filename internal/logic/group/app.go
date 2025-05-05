package group

import (
	"context"

	"gim/internal/logic/group/domain"
	"gim/internal/logic/group/repo"
	pb "gim/pkg/protocol/pb/logicpb"
)

type app struct{}

var App = new(app)

// CreateGroup 创建群组
func (*app) CreateGroup(ctx context.Context, userId uint64, in *pb.GroupCreateRequest) (uint64, error) {
	group := domain.CreateGroup(userId, in)
	err := repo.GroupRepo.Save(group)
	if err != nil {
		return 0, err
	}

	memberIDs := append([]uint64{userId}, in.MemberIds...)
	members := make([]domain.GroupUser, 0, len(memberIDs))
	for _, memberID := range memberIDs {
		memberType := pb.MemberType_GMT_MEMBER
		if memberID == userId {
			memberType = pb.MemberType_GMT_ADMIN
		}
		members = append(members, domain.GroupUser{
			GroupID:    group.ID,
			UserID:     memberID,
			MemberType: memberType,
		})
	}
	err = repo.GroupUserRepo.BatchCreate(members)
	return group.ID, err
}

// GetGroup 获取群组信息
func (*app) GetGroup(ctx context.Context, groupId uint64) (*pb.Group, error) {
	group, err := repo.GroupRepo.Get(groupId)
	if err != nil {
		return nil, err
	}

	return group.ToProto(), nil
}

// GetUserGroups 获取用户加入的群组列表
func (*app) GetUserGroups(ctx context.Context, userId uint64) ([]*pb.Group, error) {
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
func (*app) Update(ctx context.Context, userId uint64, request *pb.GroupUpdateRequest) error {
	group, err := repo.GroupRepo.Get(request.GroupId)
	if err != nil {
		return err
	}

	group.Name = request.Name
	group.AvatarUrl = request.AvatarUrl
	group.Introduction = request.Introduction
	group.Extra = request.Extra

	err = repo.GroupRepo.Save(group)
	if err != nil {
		return err
	}

	return group.PushUpdate(ctx, userId)
}

// AddMembers 添加群组成员
func (*app) AddMembers(ctx context.Context, userId, groupId uint64, userIds []uint64) error {
	group, err := repo.GroupRepo.Get(groupId)
	if err != nil {
		return err
	}
	members, err := group.AddMembers(userIds)
	if err != nil {
		return err
	}

	err = repo.GroupRepo.Save(group)
	if err != nil {
		return err
	}
	err = repo.GroupUserRepo.BatchCreate(members)
	if err != nil {
		return err
	}

	err = group.PushAddMember(ctx, userId, members)
	if err != nil {
		return err
	}
	return nil
}

// UpdateMember 更新群组用户
func (*app) UpdateMember(ctx context.Context, in *pb.UpdateMemberRequest) error {
	member, err := repo.GroupUserRepo.Get(in.GroupId, in.UserId)
	if err != nil {
		return err
	}
	member.MemberType = in.MemberType
	member.Remarks = in.Remarks
	member.Extra = in.Extra

	return repo.GroupUserRepo.Save(member)
}

// DeleteMember 删除群组成员
func (*app) DeleteMember(ctx context.Context, groupId, userId, optId uint64) error {
	err := repo.GroupUserRepo.Delete(groupId, userId)
	if err != nil {
		return err
	}

	group, err := repo.GroupRepo.Get(groupId)
	if err != nil {
		return err
	}

	return group.PushDeleteMember(ctx, optId, userId)
}

// GetMembers 获取群组成员
func (*app) GetMembers(ctx context.Context, groupId uint64) ([]*pb.GroupMember, error) {
	group, err := repo.GroupRepo.Get(groupId)
	if err != nil {
		return nil, err
	}
	return group.GetMembers(ctx)
}

// SendMessage 发送群组消息
func (*app) SendMessage(ctx context.Context, fromDeviceID, fromUserID uint64, req *pb.SendGroupMessageRequest) (uint64, error) {
	group, err := repo.GroupRepo.Get(req.GroupId)
	if err != nil {
		return 0, err
	}

	return group.SendMessage(ctx, fromDeviceID, fromUserID, req)
}

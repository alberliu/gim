package service

import (
	"context"
	"gim/internal/logic/cache"
	"gim/internal/logic/dao"
	"gim/internal/logic/model"
	"gim/pkg/pb"
	"gim/pkg/rpc"
)

type groupService struct{}

var GroupService = new(groupService)

// Get 获取群组信息
func (*groupService) Get(ctx context.Context, groupId int64) (*model.Group, error) {
	group, err := cache.GroupCache.Get(groupId)
	if err != nil {
		return nil, err
	}
	if group != nil {
		return group, nil
	}
	group, err = dao.GroupDao.Get(groupId)
	if err != nil {
		return nil, err
	}
	err = cache.GroupCache.Set(group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// Create 创建群组
func (*groupService) Create(ctx context.Context, userId int64, group model.Group, memberIds []int64) (int64, error) {
	groupId, err := dao.GroupDao.Add(group)
	if err != nil {
		return 0, err
	}

	err = SmallGroupUserService.AddUser(ctx, groupId, userId, "", "")
	if err != nil {
		return 0, err
	}
	for i := range memberIds {
		err = SmallGroupUserService.AddUser(ctx, groupId, memberIds[i], "", "")
		if err != nil {
			return 0, err
		}
	}
	return groupId, nil
}

// Update 更新群组
func (*groupService) Update(ctx context.Context, userId int64, group model.Group) error {
	err := dao.GroupDao.Update(group.Id, group.Name, group.AvatarUrl, group.Introduction, group.Extra)
	if err != nil {
		return err
	}
	err = cache.GroupCache.Del(group.Id)
	if err != nil {
		return err
	}

	userResp, err := rpc.UserIntClient.GetUser(ctx, &pb.GetUserReq{UserId: userId})
	if err != nil {
		return err
	}
	err = PushService.PushToGroup(ctx, group.Id, pb.PushCode_PC_UPDATE_GROUP, &pb.UpdateGroupPush{
		UserId:       userId,
		Nickname:     userResp.User.Nickname,
		Name:         group.Name,
		Introduction: group.Introduction,
		Extra:        group.Extra,
	}, true)
	if err != nil {
		return err
	}
	return nil
}

// GetUsers 获取群组用户
func (s *groupService) GetUsers(ctx context.Context, groupId int64) ([]*pb.GroupMember, error) {
	group, err := s.Get(ctx, groupId)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}
	if group.Type != pb.GroupType_GT_SMALL {
		return nil, nil
	}

	members, err := SmallGroupUserService.GetUsers(ctx, groupId)
	if err != nil {
		return nil, err
	}

	userIds := make(map[int64]int32, len(members))
	for i := range members {
		userIds[members[i].UserId] = 0
	}
	resp, err := rpc.UserIntClient.GetUsers(ctx, &pb.GetUsersReq{UserIds: userIds})
	if err != nil {
		return nil, err
	}

	var infos = make([]*pb.GroupMember, len(members))
	for i := range members {
		member := pb.GroupMember{
			UserId:  members[i].UserId,
			Remarks: members[i].Remarks,
			Extra:   members[i].Extra,
		}

		user, ok := resp.Users[members[i].UserId]
		if ok {
			member.Nickname = user.Nickname
			member.Sex = user.Sex
			member.AvatarUrl = user.AvatarUrl
			member.UserExtra = user.Extra
		}
		infos[i] = &member
	}

	return infos, nil
}

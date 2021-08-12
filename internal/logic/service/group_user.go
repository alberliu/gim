package service

import (
	"context"
	"gim/internal/logic/cache"
	"gim/internal/logic/dao"
	"gim/internal/logic/model"
	"gim/pkg/gerrors"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc"
)

type groupUserService struct{}

var GroupUserService = new(groupUserService)

// ListByUserId 获取用户所加入的群组
func (*groupUserService) ListByUserId(ctx context.Context, userId int64) ([]model.Group, error) {
	groups, err := dao.GroupUserDao.ListByUserId(userId)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// GetUsers 获取群组的所有用户信息
func (*groupUserService) GetUsers(ctx context.Context, groupId int64) ([]model.GroupUser, error) {
	users, err := cache.GroupUserCache.Get(groupId)
	if err != nil {
		return nil, err
	}

	if users != nil {
		return users, nil
	}

	users, err = dao.GroupUserDao.ListUser(groupId)
	if err != nil {
		return nil, err
	}

	err = cache.GroupUserCache.Set(groupId, users)
	if err != nil {
		return nil, err
	}
	return users, err
}

// AddUser 给群组添加用户
func (*groupUserService) AddUser(ctx context.Context, groupUser model.GroupUser) error {
	err := dao.GroupUserDao.Add(groupUser)
	if err != nil {
		return err
	}

	err = dao.GroupDao.UpdateUserNum(groupUser.GroupId, 1)
	if err != nil {
		return err
	}

	err = cache.GroupUserCache.Del(groupUser.GroupId)
	if err != nil {
		return err
	}

	return nil
}

// AddUsers 给群组添加用户
func (*groupUserService) AddUsers(ctx context.Context, userId, groupId int64, userIds []int64) ([]int64, error) {
	group, err := GroupService.Get(ctx, groupId)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, gerrors.ErrGroupNotExist
	}

	var existIds []int64
	var addedIds []int64

	users, err := dao.GroupUserDao.BatchGet(groupId, userIds)
	if err != nil {
		return nil, err
	}

	for i := range userIds {
		if _, ok := users[userIds[i]]; ok {
			existIds = append(existIds, userIds[i])
			continue
		}

		err = dao.GroupUserDao.Add(model.GroupUser{
			GroupId:    groupId,
			UserId:     userIds[i],
			MemberType: int(pb.MemberType_GMT_MEMBER),
		})
		if err != nil {
			return nil, err
		}

		addedIds = append(addedIds, userIds[i])
	}

	err = dao.GroupDao.UpdateUserNum(groupId, len(addedIds))
	if err != nil {
		return nil, err
	}

	err = cache.GroupUserCache.Del(groupId)
	if err != nil {
		return nil, err
	}

	var addIdMap = make(map[int64]int32, len(addedIds))
	for i := range addedIds {
		addIdMap[addedIds[i]] = 0
	}

	usersResp, err := rpc.BusinessIntClient.GetUsers(ctx, &pb.GetUsersReq{UserIds: addIdMap})
	if err != nil {
		return nil, err
	}
	var members []*pb.GroupMember
	for _, v := range usersResp.Users {
		members = append(members, &pb.GroupMember{
			UserId:    v.UserId,
			Nickname:  v.Nickname,
			Sex:       v.Sex,
			AvatarUrl: v.AvatarUrl,
			UserExtra: v.Extra,
			Remarks:   "",
			Extra:     "",
		})
	}

	userResp, err := rpc.BusinessIntClient.GetUser(ctx, &pb.GetUserReq{UserId: userId})
	if err != nil {
		return nil, err
	}

	err = PushService.PushToGroup(ctx, groupId, pb.PushCode_PC_ADD_GROUP_MEMBERS, &pb.AddGroupMembersPush{
		OptId:   userResp.User.UserId,
		OptName: userResp.User.Nickname,
		Members: members,
	}, true)
	if err != nil {
		logger.Sugar.Error(err)
	}

	return existIds, nil
}

// UpdateUser 更新群组用户
func (*groupUserService) UpdateUser(ctx context.Context, user model.GroupUser) error {
	group, err := GroupService.Get(ctx, user.GroupId)
	if err != nil {
		return err
	}

	if group == nil {
		return gerrors.ErrGroupNotExist
	}

	err = dao.GroupUserDao.Update(user)
	if err != nil {
		return err
	}

	err = cache.GroupUserCache.Del(user.GroupId)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser 删除用户群组
func (*groupUserService) DeleteUser(ctx context.Context, optId, groupId, userId int64) error {
	group, err := GroupService.Get(ctx, groupId)
	if err != nil {
		return err
	}
	if group == nil {
		return gerrors.ErrGroupNotExist
	}

	userResp, err := rpc.BusinessIntClient.GetUser(ctx, &pb.GetUserReq{UserId: optId})
	if err != nil {
		return err
	}
	err = PushService.PushToGroup(ctx, groupId, pb.PushCode_PC_REMOVE_GROUP_MEMBER, &pb.RemoveGroupMemberPush{
		OptId:         optId,
		OptName:       userResp.User.Nickname,
		DeletedUserId: userId,
	}, true)
	if err != nil {
		return err
	}

	err = dao.GroupUserDao.Delete(groupId, userId)
	if err != nil {
		return err
	}

	err = dao.GroupDao.UpdateUserNum(groupId, -1)
	if err != nil {
		return err
	}

	err = cache.GroupUserCache.Del(groupId)
	if err != nil {
		return err
	}

	return nil
}

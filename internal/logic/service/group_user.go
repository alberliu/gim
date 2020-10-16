package service

import (
	"context"
	"gim/internal/logic/cache"
	"gim/internal/logic/dao"
	"gim/pkg/gerrors"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc"
)

type groupUserService struct{}

var GroupUserService = new(groupUserService)

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

	if group.Type == pb.GroupType_GT_SMALL {
		for i := range userIds {
			member, err := dao.GroupUserDao.Get(groupId, userIds[i])
			if err != nil {
				return nil, err
			}

			if member != nil {
				existIds = append(existIds, userIds[i])
				continue
			}

			err = SmallGroupUserService.AddUser(ctx, groupId, userIds[i], "", "")
			if err != nil {
				return nil, err
			}
			addedIds = append(addedIds, userIds[i])
		}
	}
	if group.Type == pb.GroupType_GT_LARGE {
		for i := range userIds {
			isMember, err := cache.LargeGroupUserCache.IsMember(groupId, userIds[i])
			if err != nil {
				return nil, err
			}
			if isMember {
				existIds = append(existIds, userIds[i])
				continue
			}
			err = cache.LargeGroupUserCache.Set(groupId, userIds[i], "", "")
			if err != nil {
				return nil, err
			}
			addedIds = append(addedIds, userIds[i])
		}
	}

	var members []*pb.GroupMember
	for i := range addedIds {
		userResp, err := rpc.UserIntClient.GetUser(ctx, &pb.GetUserReq{UserId: addedIds[i]})
		if err != nil {
			return nil, err
		}

		members = append(members, &pb.GroupMember{
			UserId:    userResp.User.UserId,
			Nickname:  userResp.User.Nickname,
			Sex:       userResp.User.Sex,
			AvatarUrl: userResp.User.AvatarUrl,
			UserExtra: userResp.User.Extra,
			Remarks:   "",
			Extra:     "",
		})
	}

	userResp, err := rpc.UserIntClient.GetUser(ctx, &pb.GetUserReq{UserId: userId})
	if err != nil {
		return nil, err
	}

	err = PushService.PushToGroup(ctx, groupId, pb.PushCode_PC_ADD_GROUP_MEMBERS, &pb.AddGroupMembersPush{
		UserId:   userResp.User.UserId,
		Nickname: userResp.User.Nickname,
		Members:  members,
	}, true)
	if err != nil {
		logger.Sugar.Error(err)
	}

	return existIds, nil
}

// UpdateUser 更新群组用户
func (*groupUserService) UpdateUser(ctx context.Context, groupId, userId int64, label, extra string) error {
	group, err := GroupService.Get(ctx, groupId)
	if err != nil {
		return err
	}

	if group == nil {
		return gerrors.ErrGroupNotExist
	}

	if group.Type == pb.GroupType_GT_SMALL {
		err = SmallGroupUserService.Update(ctx, groupId, userId, label, extra)
		if err != nil {
			return err
		}
	}
	if group.Type == pb.GroupType_GT_LARGE {
		err = cache.LargeGroupUserCache.Set(groupId, userId, label, extra)
		if err != nil {
			return err
		}
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

	userResp, err := rpc.UserIntClient.GetUser(ctx, &pb.GetUserReq{UserId: optId})
	if err != nil {
		return err
	}
	err = PushService.PushToGroup(ctx, groupId, pb.PushCode_PC_REMOVE_GROUP_MEMBER, &pb.RemoveGroupMemberPush{
		UserId:        optId,
		Nickname:      userResp.User.Nickname,
		DeletedUserId: userId,
	}, true)
	if err != nil {
		return err
	}

	if group.Type == pb.GroupType_GT_SMALL {
		err = SmallGroupUserService.DeleteUser(ctx, groupId, userId)
		if err != nil {
			return err
		}
	}
	if group.Type == pb.GroupType_GT_LARGE {
		err = cache.LargeGroupUserCache.Del(groupId, userId)
		if err != nil {
			return err
		}
	}
	return nil
}

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

			err = SmallGroupUserService.AddUser(ctx, model.GroupUser{
				GroupId:    groupId,
				UserId:     userIds[i],
				MemberType: int(pb.MemberType_GMT_MEMBER),
			})
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
			err = cache.LargeGroupUserCache.Set(model.GroupUser{
				GroupId:    groupId,
				UserId:     userIds[i],
				MemberType: int(pb.MemberType_GMT_MEMBER),
			})
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

	if group.Type == pb.GroupType_GT_SMALL {
		err = SmallGroupUserService.Update(ctx, user)
		if err != nil {
			return err
		}
	}
	if group.Type == pb.GroupType_GT_LARGE {
		err = cache.LargeGroupUserCache.Set(user)
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
		OptId:         optId,
		OptName:       userResp.User.Nickname,
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

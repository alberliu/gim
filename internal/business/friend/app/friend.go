package app

import (
	"context"
	"errors"
	"time"

	"gim/internal/business/friend/domain"
	"gim/internal/business/friend/repo"
	userapp "gim/internal/business/user/app"
	"gim/pkg/gerrors"
	pb "gim/pkg/protocol/pb/businesspb"
	"gim/pkg/protocol/pb/connectpb"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/rpc"
)

var FriendApp = new(friendApp)

type friendApp struct{}

// List 获取好友列表
func (s *friendApp) List(ctx context.Context, userID uint64) ([]*pb.Friend, error) {
	friends, err := repo.FriendRepo.List(ctx, userID, domain.FriendStatusAgree)
	if err != nil {
		return nil, err
	}

	userIDs := make([]uint64, 0, len(friends))
	for i := range friends {
		userIDs = append(userIDs, friends[i].UserID)
	}
	users, err := userapp.UserApp.GetUsers(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	var infos = make([]*pb.Friend, len(friends))
	for i := range friends {
		friend := pb.Friend{
			UserId:  friends[i].FriendID,
			Remarks: friends[i].Remarks,
			Extra:   friends[i].Extra,
		}

		user, ok := users[friends[i].FriendID]
		if ok {
			friend.Nickname = user.Nickname
			friend.Sex = user.Sex
			friend.AvatarUrl = user.AvatarUrl
			friend.UserExtra = user.Extra
		}
		infos[i] = &friend
	}

	return infos, nil
}

// AddFriend 添加好友
func (*friendApp) AddFriend(ctx context.Context, userId, friendId uint64, remarks, description string) error {
	friend, err := repo.FriendRepo.Get(ctx, userId, friendId)
	if err != nil && !errors.Is(err, gerrors.ErrFriendNotFound) {
		return err
	}
	if err == nil {
		if friend.Status == domain.FriendStatusApply {
			return nil
		}
		if friend.Status == domain.FriendStatusAgree {
			return gerrors.ErrAlreadyIsFriend
		}
	}

	err = repo.FriendRepo.Create(ctx, &domain.Friend{
		UserID:   userId,
		FriendID: friendId,
		Remarks:  remarks,
		Status:   domain.FriendStatusApply,
	})
	if err != nil {
		return err
	}

	user, err := userapp.UserApp.Get(ctx, userId)
	if err != nil {
		return err
	}

	_, err = rpc.PushToUsers(ctx, rpc.PushRequest{
		UserIDs: []uint64{friendId},
		Command: connectpb.MessageCommand(pb.Command_ADD_FRIEND),
		Message: &pb.AddFriendPush{
			FriendId:    userId,
			Nickname:    user.Nickname,
			AvatarUrl:   user.AvatarUrl,
			Description: description,
		},
		IsPersist: true,
	})
	return err
}

// AgreeAddFriend 同意添加好友
func (*friendApp) AgreeAddFriend(ctx context.Context, userId, friendId uint64, remarks string) error {
	friend, err := repo.FriendRepo.Get(ctx, friendId, userId)
	if err != nil {
		return err
	}
	if friend.Status == domain.FriendStatusAgree {
		return nil
	}
	friend.Status = domain.FriendStatusAgree
	err = repo.FriendRepo.Save(ctx, friend)
	if err != nil {
		return err
	}

	err = repo.FriendRepo.Save(ctx, &domain.Friend{
		UserID:   userId,
		FriendID: friendId,
		Remarks:  remarks,
		Status:   domain.FriendStatusAgree,
	})
	if err != nil {
		return err
	}

	user, err := userapp.UserApp.Get(ctx, userId)
	if err != nil {
		return err
	}

	_, err = rpc.PushToUsers(ctx, rpc.PushRequest{
		UserIDs: []uint64{friendId},
		Command: connectpb.MessageCommand(pb.Command_AGREE_ADD_FRIEND),
		Message: &pb.AgreeAddFriendPush{
			FriendId:  userId,
			Nickname:  user.Nickname,
			AvatarUrl: user.AvatarUrl,
		},
		IsPersist: true,
	})
	return err
}

// SetFriend 设置好友信息
func (*friendApp) SetFriend(ctx context.Context, userId uint64, req *pb.FriendSetRequest) error {
	friend, err := repo.FriendRepo.Get(ctx, userId, req.FriendId)
	if err != nil {
		return err
	}

	friend.Remarks = req.Remarks
	friend.Extra = req.Extra
	friend.UpdatedAt = time.Now()

	return repo.FriendRepo.Save(ctx, friend)
}

// SendToFriend 消息发送至好友
func (*friendApp) SendToFriend(ctx context.Context, fromDeviceID, fromUserID uint64, req *pb.SendFriendMessageRequest) (uint64, error) {
	user, err := userapp.UserApp.Get(ctx, fromUserID)
	if err != nil {
		return 0, err
	}

	// 发给发送者
	push := &logicpb.UserMessagePush{
		FromUser: &logicpb.User{
			UserId:    fromUserID,
			DeviceId:  fromDeviceID,
			Nickname:  user.Nickname,
			AvatarUrl: user.AvatarUrl,
			Extra:     user.Extra,
		},
		ToUserId: req.UserId,
		Content:  req.Content,
	}

	pushReply, err := rpc.PushToUsers(ctx, rpc.PushRequest{
		UserIDs:   []uint64{fromUserID, req.UserId},
		Command:   connectpb.MessageCommand_MC_USER_MESSAGE,
		Message:   push,
		IsPersist: true,
	})
	if err != nil {
		return 0, err
	}
	return pushReply.MessageId, nil
}

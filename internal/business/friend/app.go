package friend

import (
	"context"
	"errors"
	"time"

	userapp "gim/internal/business/user/app"
	"gim/internal/logic/message"
	"gim/pkg/gerrors"
	pb "gim/pkg/protocol/pb/businesspb"
	"gim/pkg/protocol/pb/connectpb"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/rpc"
)

type app struct{}

var App = new(app)

// List 获取好友列表
func (s *app) List(ctx context.Context, userID uint64) ([]*pb.Friend, error) {
	friends, err := Repo.List(userID, FriendStatusAgree)
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
func (*app) AddFriend(ctx context.Context, userId, friendId uint64, remarks, description string) error {
	friend, err := Repo.Get(userId, friendId)
	if err != nil && !errors.Is(err, gerrors.ErrFriendNotFound) {
		return err
	}
	if err == nil {
		if friend.Status == FriendStatusApply {
			return nil
		}
		if friend.Status == FriendStatusAgree {
			return gerrors.ErrAlreadyIsFriend
		}
	}

	err = Repo.Create(&Friend{
		UserID:   userId,
		FriendID: friendId,
		Remarks:  remarks,
		Status:   FriendStatusApply,
	})
	if err != nil {
		return err
	}

	user, err := userapp.UserApp.Get(ctx, userId)
	if err != nil {
		return err
	}

	_, err = message.App.PushToUsersWithAny(ctx, []uint64{friendId}, connectpb.Command_ADD_FRIEND, &pb.AddFriendPush{
		FriendId:    userId,
		Nickname:    user.Nickname,
		AvatarUrl:   user.AvatarUrl,
		Description: description,
	}, true)
	return err
}

// AgreeAddFriend 同意添加好友
func (*app) AgreeAddFriend(ctx context.Context, userId, friendId uint64, remarks string) error {
	friend, err := Repo.Get(friendId, userId)
	if err != nil {
		return err
	}
	if friend.Status == FriendStatusAgree {
		return nil
	}
	friend.Status = FriendStatusAgree
	err = Repo.Save(friend)
	if err != nil {
		return err
	}

	err = Repo.Save(&Friend{
		UserID:   userId,
		FriendID: friendId,
		Remarks:  remarks,
		Status:   FriendStatusAgree,
	})
	if err != nil {
		return err
	}

	user, err := userapp.UserApp.Get(ctx, userId)
	if err != nil {
		return err
	}

	_, err = message.App.PushToUsersWithAny(ctx, []uint64{friendId}, connectpb.Command_AGREE_ADD_FRIEND, &pb.AgreeAddFriendPush{
		FriendId:  userId,
		Nickname:  user.Nickname,
		AvatarUrl: user.AvatarUrl,
	}, true)
	return err
}

// SetFriend 设置好友信息
func (*app) SetFriend(ctx context.Context, userId uint64, req *pb.FriendSetRequest) error {
	friend, err := Repo.Get(userId, req.FriendId)
	if err != nil {
		return err
	}

	friend.Remarks = req.Remarks
	friend.Extra = req.Extra
	friend.UpdatedAt = time.Now()

	return Repo.Save(friend)
}

// SendToFriend 消息发送至好友
func (*app) SendToFriend(ctx context.Context, fromDeviceID, fromUserID uint64, req *pb.SendFriendMessageRequest) (uint64, error) {
	user, err := rpc.GetUser(fromDeviceID, fromUserID)
	if err != nil {
		return 0, err
	}

	// 发给发送者
	push := logicpb.UserMessagePush{
		FromUser: user,
		ToUserId: req.UserId,
		Content:  req.Content,
	}
	userIDs := []uint64{fromUserID, req.UserId}
	return message.App.PushToUsersWithAny(ctx, userIDs, connectpb.Command_USER_MESSAGE, &push, true)
}

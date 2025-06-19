package friend

import (
	"context"
	"errors"
	"time"

	"google.golang.org/protobuf/proto"

	"gim/internal/logic/message"
	"gim/pkg/gerrors"
	pb "gim/pkg/protocol/pb/logicpb"
	"gim/pkg/protocol/pb/userpb"
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

	userIDs := make(map[uint64]int32, len(friends))
	for i := range friends {
		userIDs[friends[i].FriendID] = 0
	}
	reply, err := rpc.GetUserIntClient().GetUsers(ctx, &userpb.GetUsersRequest{UserIds: userIDs})
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

		user, ok := reply.Users[friends[i].FriendID]
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

	reply, err := rpc.GetUserIntClient().GetUser(ctx, &userpb.GetUserRequest{UserId: userId})
	if err != nil {
		return err
	}

	_, err = message.App.PushToUser(ctx, []uint64{friendId}, pb.PushCode_PC_ADD_FRIEND, &pb.AddFriendPush{
		FriendId:    userId,
		Nickname:    reply.User.Nickname,
		AvatarUrl:   reply.User.AvatarUrl,
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

	reply, err := rpc.GetUserIntClient().GetUser(ctx, &userpb.GetUserRequest{UserId: userId})
	if err != nil {
		return err
	}

	_, err = message.App.PushToUser(ctx, []uint64{friendId}, pb.PushCode_PC_AGREE_ADD_FRIEND, &pb.AgreeAddFriendPush{
		FriendId:  userId,
		Nickname:  reply.User.Nickname,
		AvatarUrl: reply.User.AvatarUrl,
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
	sender, err := rpc.GetSender(fromDeviceID, fromUserID)
	if err != nil {
		return 0, err
	}

	// 发给发送者
	push := pb.UserMessagePush{
		Sender:  sender,
		Content: req.Content,
	}
	bytes, err := proto.Marshal(&push)
	if err != nil {
		return 0, err
	}

	msg := &pb.Message{
		Code:    pb.PushCode_PC_USER_MESSAGE,
		Content: bytes,
	}

	userIDs := []uint64{fromUserID, req.UserId}
	return message.App.SendToUsers(ctx, userIDs, msg, true)
}

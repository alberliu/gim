package app

import (
	"context"
	frienddomain "gim/internal/logic/domain/friend"
	"gim/pkg/pb"
	"time"
)

type friendApp struct{}

var FriendApp = new(friendApp)

// List 获取好友列表
func (s *friendApp) List(ctx context.Context, userId int64) ([]*pb.Friend, error) {
	return frienddomain.FriendService.List(ctx, userId)
}

// AddFriend 添加好友
func (*friendApp) AddFriend(ctx context.Context, userId, friendId int64, remarks, description string) error {
	return frienddomain.FriendService.AddFriend(ctx, userId, friendId, remarks, description)
}

// AgreeAddFriend 同意添加好友
func (*friendApp) AgreeAddFriend(ctx context.Context, userId, friendId int64, remarks string) error {
	return frienddomain.FriendService.AgreeAddFriend(ctx, userId, friendId, remarks)
}

// SetFriend 设置好友信息
func (*friendApp) SetFriend(ctx context.Context, userId int64, req *pb.SetFriendReq) error {
	friend, err := frienddomain.FriendRepo.Get(userId, req.FriendId)
	if err != nil {
		return err
	}
	if friend == nil {
		return nil
	}

	friend.Remarks = req.Remarks
	friend.Extra = req.Extra
	friend.UpdateTime = time.Now()

	err = frienddomain.FriendRepo.Save(friend)
	if err != nil {
		return err
	}
	return nil
}

// SendToFriend 消息发送至好友
func (*friendApp) SendToFriend(ctx context.Context, sender *pb.Sender, req *pb.SendMessageReq) (int64, error) {
	return frienddomain.FriendService.SendToFriend(ctx, sender, req)
}

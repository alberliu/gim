package app

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	userapp "gim/internal/business/user/app"
	pb "gim/pkg/gen/proto/businesspb"
	"gim/pkg/gen/proto/connectpb"
	"gim/pkg/gen/proto/logicpb"
	"gim/pkg/rpc"
)

var MessageApp = new(messageApp)

type messageApp struct{}

func (*messageApp) SendToFriend(ctx context.Context, fromDeviceID, fromUserID uint64, req *pb.SendFriendMessageRequest) (
	*pb.SendFriendMessageReply, error) {
	user, err := userapp.UserApp.Get(ctx, fromUserID)
	if err != nil {
		return nil, err
	}

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

	reply, err := rpc.PushToUsers(ctx, rpc.PushRequest{
		FromUserID: fromUserID,
		UserIDs:    []uint64{fromUserID, req.UserId},
		Command:    connectpb.MessageCommand_MC_USER_MESSAGE,
		Message:    push,
		IsPersist:  true,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SendFriendMessageReply{
		MessageId:   reply.MessageId,
		FromUserSeq: reply.FromUserSeq,
	}, nil
}

func (*messageApp) SendToGroup(ctx context.Context, fromDeviceID, fromUserID uint64, req *pb.SendGroupMessageRequest) (*pb.SendGroupMessageReply, error) {
	user, err := userapp.UserApp.Get(ctx, fromUserID)
	if err != nil {
		return nil, err
	}

	push := &logicpb.GroupMessagePush{
		FromUser: &logicpb.User{
			UserId:    fromUserID,
			DeviceId:  fromDeviceID,
			Nickname:  user.Nickname,
			AvatarUrl: user.AvatarUrl,
			Extra:     user.Extra,
		},
		GroupId: req.GroupId,
		Content: req.Content,
	}

	content, err := proto.Marshal(push)
	if err != nil {
		return nil, err
	}

	reply, err := rpc.GetGroupIntClient().Push(ctx, &logicpb.GroupPushRequest{
		FromUserId: fromUserID,
		GroupId:    req.GroupId,
		Command:    connectpb.MessageCommand_MC_GROUP_MESSAGE,
		Content:    content,
		IsPersist:  true,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SendGroupMessageReply{
		MessageId:   reply.MessageId,
		FromUserSeq: reply.FromUserSeq,
	}, nil
}

func (*messageApp) SendToRoom(ctx context.Context, req *pb.SendRoomMessageRequest) error {
	_, err := rpc.GetRoomIntClient().PushRoom(ctx, &logicpb.PushRoomRequest{
		RoomId:   req.RoomId,
		Command:  connectpb.MessageCommand_MC_ROOM_MESSAGE,
		Content:  req.Content,
		SendTime: time.Now().Unix(),
	})
	return err
}

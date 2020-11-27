package service

import (
	"context"
	"gim/internal/logic/cache"
	"gim/internal/logic/dao"
	"gim/internal/logic/model"
	"gim/pkg/gerrors"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc"
	"gim/pkg/util"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

const MessageLimit = 50 // 最大消息同步数量

const MaxSyncBufLen = 65536 // 最大字节数组长度

type messageService struct{}

var MessageService = new(messageService)

// Add 添加消息
func (*messageService) Add(ctx context.Context, message model.Message) error {
	return dao.MessageDao.Add("message", message)
}

// Sync 消息同步
func (*messageService) Sync(ctx context.Context, userId, seq int64) (*pb.SyncResp, error) {
	messages, hasMore, err := MessageService.ListByUserIdAndSeq(ctx, userId, seq)
	if err != nil {
		return nil, err
	}
	pbMessages := model.MessagesToPB(messages)
	length := len(pbMessages)

	resp := &pb.SyncResp{Messages: pbMessages, HasMore: hasMore}
	bytes, err := proto.Marshal(resp)
	if err != nil {
		return nil, err
	}

	// 如果字节数组大于一个包的长度，需要减少字节数组
	for len(bytes) > MaxSyncBufLen {
		length = length * 2 / 3
		resp = &pb.SyncResp{Messages: pbMessages[0:length], HasMore: true}
		bytes, err = proto.Marshal(resp)
		if err != nil {
			return nil, err
		}
	}

	var userIds = make(map[int64]int32, len(resp.Messages))
	for i := range resp.Messages {
		if resp.Messages[i].Sender.SenderType == pb.SenderType_ST_USER {
			userIds[resp.Messages[i].Sender.SenderId] = 0
		}
	}
	usersResp, err := rpc.UserIntClient.GetUsers(ctx, &pb.GetUsersReq{UserIds: userIds})
	if err != nil {
		return nil, err
	}
	for i := range resp.Messages {
		if resp.Messages[i].Sender.SenderType == pb.SenderType_ST_USER {
			user, ok := usersResp.Users[resp.Messages[i].Sender.SenderId]
			if ok {
				resp.Messages[i].Sender.Nickname = user.Nickname
				resp.Messages[i].Sender.AvatarUrl = user.AvatarUrl
				resp.Messages[i].Sender.Extra = user.Extra
			} else {
				logger.Logger.Warn("get user failed", zap.Int64("user_id", resp.Messages[i].Sender.SenderId))
			}
		}
	}

	return resp, nil
}

// ListByUserIdAndSeq 查询消息
func (*messageService) ListByUserIdAndSeq(ctx context.Context, userId, seq int64) ([]model.Message, bool, error) {
	var err error
	if seq == 0 {
		seq, err = DeviceAckService.GetMaxByUserId(ctx, userId)
		if err != nil {
			return nil, false, err
		}
	}
	return dao.MessageDao.ListBySeq("message", model.MessageObjectTypeUser, userId, seq, MessageLimit)
}

// Send 消息发送
func (s *messageService) Send(ctx context.Context, sender model.Sender, req pb.SendMessageReq) (int64, error) {
	// 如果发送者是用户，需要补充用户的信息
	if sender.SenderType == pb.SenderType_ST_USER {
		user, err := rpc.UserIntClient.GetUser(ctx, &pb.GetUserReq{UserId: sender.SenderId})
		if err != nil {
			return 0, err
		}
		if user.User == nil {
			return 0, gerrors.ErrUserNotFound
		}
		sender.AvatarUrl = user.User.AvatarUrl
		sender.Nickname = user.User.Nickname
		sender.Extra = user.User.Extra
	}

	switch req.ReceiverType {
	// 消息接收者为用户
	case pb.ReceiverType_RT_USER:
		// 发送者为用户
		if sender.SenderType == pb.SenderType_ST_USER {
			return MessageService.SendToFriend(ctx, sender, req)
		} else {
			return MessageService.SendToUser(ctx, sender, req.ReceiverId, req)
		}
	// 消息接收者是小群组
	case pb.ReceiverType_RT_SMALL_GROUP:
		return MessageService.SendToGroup(ctx, sender, req)
	// 消息接收者是大群组
	case pb.ReceiverType_RT_LARGE_GROUP:
		return MessageService.SendToLargeGroup(ctx, sender, req)
	}
	return 0, nil
}

// SendToUser 消息发送至用户
func (*messageService) SendToFriend(ctx context.Context, sender model.Sender, req pb.SendMessageReq) (int64, error) {
	// 发给发送者
	seq, err := MessageService.SendToUser(ctx, sender, sender.SenderId, req)
	if err != nil {
		return 0, err
	}

	// 用户需要增加自己的已经同步的序列号
	err = cache.DeviceACKCache.Set(sender.SenderId, sender.DeviceId, seq)
	if err != nil {
		return 0, err
	}

	// 发给接收者
	_, err = MessageService.SendToUser(ctx, sender, req.ReceiverId, req)
	if err != nil {
		return 0, err
	}

	return seq, nil
}

// SendToGroup 消息发送至群组（使用写扩散）
func (*messageService) SendToGroup(ctx context.Context, sender model.Sender, req pb.SendMessageReq) (int64, error) {
	users, err := SmallGroupUserService.GetUsers(ctx, req.ReceiverId)
	if err != nil {
		return 0, err
	}

	if sender.SenderType == pb.SenderType_ST_USER && !IsInGroup(users, sender.SenderId) {
		logger.Sugar.Error(ctx, sender.SenderId, req.ReceiverId, "不在群组内")
		return 0, gerrors.ErrNotInGroup
	}

	var userSeq int64
	// 将消息发送给群组用户，使用写扩散
	for _, user := range users {
		seq, err := MessageService.SendToUser(ctx, sender, user.UserId, req)
		if err != nil {
			return 0, err
		}
		if user.UserId == sender.SenderId {
			userSeq = seq
		}
	}

	if sender.SenderType == pb.SenderType_ST_USER {
		// 用户需要增加自己的已经同步的序列号
		err = cache.DeviceACKCache.Set(sender.SenderId, sender.DeviceId, userSeq)
		if err != nil {
			return 0, err
		}
	}

	return userSeq, nil
}

func IsInGroup(users []model.GroupUser, userId int64) bool {
	for i := range users {
		if users[i].UserId == userId {
			return true
		}
	}
	return false
}

// SendToLargeGroup 消息发送至大群组（读扩散）
func (*messageService) SendToLargeGroup(ctx context.Context, sender model.Sender, req pb.SendMessageReq) (int64, error) {
	users, err := cache.LargeGroupUserCache.Members(req.ReceiverId)
	if err != nil {
		return 0, err
	}

	isMember, err := cache.LargeGroupUserCache.IsMember(req.ReceiverId, sender.SenderId)
	if err != nil {
		return 0, err
	}

	if sender.SenderType == pb.SenderType_ST_USER && !isMember {
		logger.Logger.Warn("not int group", zap.Int64("group_id", req.ReceiverId), zap.Int64("user_id", sender.SenderId))
		return 0, gerrors.ErrNotInGroup
	}

	var seq int64 = 0
	if req.IsPersist {
		seq, err = SeqService.GetGroupNext(ctx, req.ReceiverId)
		if err != nil {
			return 0, err
		}
	}

	go func() {
		defer util.RecoverPanic()

		if req.IsPersist {
			message := model.Message{
				ObjectType:   model.MessageObjectTypeGroup,
				ObjectId:     req.ReceiverId,
				RequestId:    grpclib.GetCtxRequstId(ctx),
				SenderType:   int32(sender.SenderType),
				SenderId:     sender.SenderId,
				ReceiverType: int32(req.ReceiverType),
				ReceiverId:   req.ReceiverId,
				ToUserIds:    model.FormatUserIds(req.ToUserIds),
				Type:         int(req.MessageType),
				Content:      req.MessageContent,
				Seq:          seq,
				SendTime:     util.UnunixMilliTime(req.SendTime),
				Status:       int32(pb.MessageStatus_MS_NORMAL),
			}
			err = MessageService.Add(ctx, message)
			if err != nil {
				logger.Sugar.Error(err)
				return
			}
		}

		// 将消息发送给群组用户，使用读扩散
		for i := range users {
			err = MessageService.SendToLargeGroupUser(ctx, sender, users[i].UserId, seq, req)
			if err != nil {
				return
			}
		}
	}()

	return seq, nil
}

// SendToUser 将消息发送给用户
func (*messageService) SendToUser(ctx context.Context, sender model.Sender, toUserId int64, req pb.SendMessageReq) (int64, error) {
	logger.Logger.Debug("SendToUser",
		zap.Int64("request_id", grpclib.GetCtxRequstId(ctx)),
		zap.Int64("to_user_id", toUserId))

	var (
		seq int64 = 0
		err error
	)

	if req.IsPersist {
		seq, err = SeqService.GetUserNext(ctx, toUserId)
		if err != nil {
			return 0, err
		}
	}

	go func() {
		defer util.RecoverPanic()

		if req.IsPersist {
			selfMessage := model.Message{
				ObjectType:   model.MessageObjectTypeUser,
				ObjectId:     toUserId,
				RequestId:    grpclib.GetCtxRequstId(ctx),
				SenderType:   int32(sender.SenderType),
				SenderId:     sender.SenderId,
				ReceiverType: int32(req.ReceiverType),
				ReceiverId:   req.ReceiverId,
				ToUserIds:    model.FormatUserIds(req.ToUserIds),
				Type:         int(req.MessageType),
				Content:      req.MessageContent,
				Seq:          seq,
				SendTime:     util.UnunixMilliTime(req.SendTime),
				Status:       int32(pb.MessageStatus_MS_NORMAL),
			}
			err = MessageService.Add(ctx, selfMessage)
			if err != nil {
				logger.Sugar.Error(err)
				return
			}
		}

		message := pb.Message{
			Sender: &pb.Sender{
				SenderType: sender.SenderType,
				SenderId:   sender.SenderId,
				AvatarUrl:  sender.AvatarUrl,
				Nickname:   sender.Nickname,
				Extra:      sender.Extra,
			},
			ReceiverType:   req.ReceiverType,
			ReceiverId:     req.ReceiverId,
			ToUserIds:      req.ToUserIds,
			MessageType:    req.MessageType,
			MessageContent: req.MessageContent,
			Seq:            seq,
			SendTime:       req.SendTime,
			Status:         pb.MessageStatus_MS_NORMAL,
		}

		// 查询用户在线设备
		devices, err := DeviceService.ListOnlineByUserId(ctx, toUserId)
		if err != nil {
			logger.Sugar.Error(err)
			return
		}

		for i := range devices {
			// 消息不需要投递给发送消息的设备
			if sender.DeviceId == devices[i].Id {
				continue
			}

			err = MessageService.SendToDevice(context.TODO(), devices[i], message)
			if err != nil {
				logger.Sugar.Error(err, zap.Any("context canceled", devices[i]))
				return
			}
		}
	}()

	return seq, nil
}

// SendToLargeGroupUser 发送消息给大群组用户
func (*messageService) SendToLargeGroupUser(ctx context.Context, sender model.Sender, toUserId int64, roomSeq int64, req pb.SendMessageReq) error {
	logger.Logger.Debug("SendToLargeGroupUser",
		zap.Int64("request_id", grpclib.GetCtxRequstId(ctx)),
		zap.Int64("to_user_id", toUserId))

	message := pb.Message{
		Sender: &pb.Sender{
			SenderType: sender.SenderType,
			SenderId:   sender.SenderId,
			AvatarUrl:  sender.AvatarUrl,
			Nickname:   sender.Nickname,
			Extra:      sender.Extra,
		},
		ReceiverType:   req.ReceiverType,
		ReceiverId:     req.ReceiverId,
		ToUserIds:      req.ToUserIds,
		MessageType:    req.MessageType,
		MessageContent: req.MessageContent,
		Seq:            roomSeq,
		SendTime:       req.SendTime,
		Status:         pb.MessageStatus_MS_NORMAL,
	}

	// 查询用户在线设备
	devices, err := DeviceService.ListOnlineByUserId(ctx, toUserId)
	if err != nil {
		return err
	}

	for i := range devices {
		// 消息不需要投递给发送消息的设备
		if sender.DeviceId == devices[i].Id {
			continue
		}

		err = MessageService.SendToDevice(ctx, devices[i], message)
		if err != nil {
			return err
		}
	}
	return nil
}

// SendToDevice 将消息发送给设备
func (*messageService) SendToDevice(ctx context.Context, device model.Device, message pb.Message) error {
	if device.Status == model.DeviceOnLine {
		messageSend := pb.MessageSend{Message: &message}
		_, err := rpc.ConnectIntClient.DeliverMessage(grpclib.ContextWithAddr(ctx, device.ConnAddr), &pb.DeliverMessageReq{
			DeviceId:    device.Id,
			Fd:          device.ConnFd,
			MessageSend: &messageSend,
		})
		if err != nil {
			return err
		}
	}

	// todo 其他推送厂商
	return nil
}

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

	"go.uber.org/zap"
)

type messageService struct{}

var MessageService = new(messageService)

// Add 添加消息
func (*messageService) Add(ctx context.Context, message model.Message) error {
	return dao.MessageDao.Add("message", message)
}

// ListByUserIdAndSeq 查询消息
func (*messageService) ListByUserIdAndSeq(ctx context.Context, userId, seq int64) ([]model.Message, error) {
	var err error
	if seq == 0 {
		seq, err = DeviceAckService.GetMaxByUserId(ctx, userId)
		if err != nil {
			return nil, err
		}
	}
	messages, err := dao.MessageDao.ListBySeq("message", model.MessageObjectTypeUser, userId, seq)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// Send 消息发送
func (s *messageService) Send(ctx context.Context, sender model.Sender, req pb.SendMessageReq) (int64, error) {
	switch req.ReceiverType {
	case pb.ReceiverType_RT_USER:
		if sender.SenderType == pb.SenderType_ST_USER {
			return MessageService.SendToFriend(ctx, sender, req)
		} else {
			return MessageService.SendToUser(ctx, sender, req.ReceiverId, 0, req)
		}
	case pb.ReceiverType_RT_NORMAL_GROUP:
		return MessageService.SendToGroup(ctx, sender, req)

	case pb.ReceiverType_RT_LARGE_GROUP:
		return MessageService.SendToLargeGroup(ctx, sender, req)
	}
	return 0, nil
}

// SendToUser 消息发送至用户
func (*messageService) SendToFriend(ctx context.Context, sender model.Sender, req pb.SendMessageReq) (int64, error) {
	// 发给发送者
	seq, err := MessageService.SendToUser(ctx, sender, sender.SenderId, 0, req)
	if err != nil {
		return 0, err
	}

	// 发给接收者
	_, err = MessageService.SendToUser(ctx, sender, req.ReceiverId, 0, req)
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
		seq, err := MessageService.SendToUser(ctx, sender, user.UserId, 0, req)
		if err != nil {
			return 0, err
		}
		if user.UserId == sender.SenderId {
			userSeq = seq
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
		logger.Logger.Error("not int group", zap.Int64("group_id", req.ReceiverId), zap.Int64("user_id", sender.SenderId))
		return 0, gerrors.ErrNotInGroup
	}

	var seq int64 = 0
	if req.IsPersist {
		seq, err = SeqService.GetGroupNext(ctx, req.ReceiverId)
		if err != nil {
			return 0, err
		}
		message := model.Message{
			ObjectType:     model.MessageObjectTypeGroup,
			ObjectId:       req.ReceiverId,
			RequestId:      grpclib.GetCtxRequstId(ctx),
			SenderType:     int32(sender.SenderType),
			SenderId:       sender.SenderId,
			SenderDeviceId: sender.DeviceId,
			ReceiverType:   int32(req.ReceiverType),
			ReceiverId:     req.ReceiverId,
			ToUserIds:      model.FormatUserIds(req.ToUserIds),
			Type:           int(req.MessageType),
			Content:        req.MessageContent,
			Seq:            seq,
			SendTime:       util.UnunixMilliTime(req.SendTime),
			Status:         int32(pb.MessageStatus_MS_NORMAL),
		}
		err = MessageService.Add(ctx, message)
		if err != nil {
			return 0, err
		}
	}

	// 将消息发送给群组用户，使用读扩散
	req.IsPersist = false
	for i := range users {
		_, err = MessageService.SendToUser(ctx, sender, users[i].UserId, seq, req)
		if err != nil {
			return 0, err
		}
	}
	return seq, nil
}

// StoreAndSendToUser 将消息持久化到数据库,并且消息发送至用户
func (*messageService) SendToUser(ctx context.Context, sender model.Sender, toUserId int64, roomSeq int64, req pb.SendMessageReq) (int64, error) {
	logger.Logger.Debug("message_store_send_to_user",
		zap.Int64("request_id", grpclib.GetCtxRequstId(ctx)),
		zap.Int64("to_user_id", toUserId))

	var (
		seq = roomSeq
		err error
	)
	if req.IsPersist {
		seq, err = SeqService.GetUserNext(ctx, toUserId)
		if err != nil {
			return 0, err
		}
		selfMessage := model.Message{
			ObjectType:     model.MessageObjectTypeUser,
			ObjectId:       toUserId,
			RequestId:      grpclib.GetCtxRequstId(ctx),
			SenderType:     int32(sender.SenderType),
			SenderId:       sender.SenderId,
			SenderDeviceId: sender.DeviceId,
			ReceiverType:   int32(req.ReceiverType),
			ReceiverId:     req.ReceiverId,
			ToUserIds:      model.FormatUserIds(req.ToUserIds),
			Type:           int(req.MessageType),
			Content:        req.MessageContent,
			Seq:            seq,
			SendTime:       util.UnunixMilliTime(req.SendTime),
			Status:         int32(pb.MessageStatus_MS_NORMAL),
		}
		err = MessageService.Add(ctx, selfMessage)
		if err != nil {
			return 0, err
		}
	}

	messageItem := pb.MessageItem{
		RequestId:      grpclib.GetCtxRequstId(ctx),
		SenderType:     sender.SenderType,
		SenderId:       sender.SenderId,
		SenderDeviceId: sender.DeviceId,
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
		return 0, err
	}

	for i := range devices {
		// 消息不需要投递给发送消息的设备
		if sender.DeviceId == devices[i].Id {
			continue
		}

		err = MessageService.SendToDevice(ctx, devices[i], messageItem)
		if err != nil {
			return 0, err
		}
	}
	return seq, nil
}

// SendToDevice 将消息发送给设备
func (*messageService) SendToDevice(ctx context.Context, device model.Device, msgItem pb.MessageItem) error {
	if device.Status == model.DeviceOnLine {
		message := pb.Message{Message: &msgItem}
		_, err := rpc.ConnectIntClient.DeliverMessage(grpclib.ContextWithAddr(ctx, device.ConnAddr), &pb.DeliverMessageReq{
			DeviceId: device.Id,
			Fd:       device.ConnFd,
			Message:  &message,
		})
		if err != nil {
			return err
		}
	}

	// todo 其他推送厂商
	return nil
}

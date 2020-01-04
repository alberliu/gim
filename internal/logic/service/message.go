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
	"gim/pkg/rpc_cli"
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
func (*messageService) ListByUserIdAndSeq(ctx context.Context, appId, userId, seq int64) ([]model.Message, error) {
	var err error
	if seq == 0 {
		seq, err = DeviceAckService.GetMaxByUserId(ctx, appId, userId)
		if err != nil {
			return nil, err
		}
	}
	messages, err := dao.MessageDao.ListBySeq("message", appId, model.MessageObjectTypeUser, userId, seq)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// Send 消息发送
func (s *messageService) Send(ctx context.Context, sender model.Sender, req pb.SendMessageReq) error {
	switch req.ReceiverType {
	case pb.ReceiverType_RT_USER:
		if sender.SenderType == pb.SenderType_ST_USER {
			err := MessageService.SendToFriend(ctx, sender, req)
			if err != nil {
				return err
			}
		} else {
			err := MessageService.SendToUser(ctx, sender, req.ReceiverId, 0, req)
			if err != nil {
				return err
			}
		}
	case pb.ReceiverType_RT_NORMAL_GROUP:
		err := MessageService.SendToGroup(ctx, sender, req)
		if err != nil {
			return err
		}
	case pb.ReceiverType_RT_LARGE_GROUP:
		err := MessageService.SendToChatRoom(ctx, sender, req)
		if err != nil {
			return err
		}
	}

	return nil
}

// SendToUser 消息发送至用户
func (*messageService) SendToFriend(ctx context.Context, sender model.Sender, req pb.SendMessageReq) error {
	// 发给发送者
	err := MessageService.SendToUser(ctx, sender, sender.SenderId, 0, req)
	if err != nil {
		return err
	}

	// 发给接收者
	err = MessageService.SendToUser(ctx, sender, req.ReceiverId, 0, req)
	if err != nil {
		return err
	}

	return nil
}

// SendToGroup 消息发送至群组（使用写扩散）
func (*messageService) SendToGroup(ctx context.Context, sender model.Sender, req pb.SendMessageReq) error {
	users, err := GroupUserService.GetUsers(ctx, sender.AppId, req.ReceiverId)
	if err != nil {
		return err
	}

	if sender.SenderType == pb.SenderType_ST_USER && !IsInGroup(users, sender.SenderId) {
		logger.Sugar.Error(ctx, sender.SenderId, req.ReceiverId, "不在群组内")
		return gerrors.ErrNotInGroup
	}

	// 将消息发送给群组用户，使用写扩散
	for _, user := range users {
		err = MessageService.SendToUser(ctx, sender, user.UserId, 0, req)
		if err != nil {
			return err
		}
	}
	return nil
}

func IsInGroup(users []model.GroupUser, userId int64) bool {
	for i := range users {
		if users[i].UserId == userId {
			return true
		}
	}
	return false
}

// SendToChatRoom 消息发送至聊天室（读扩散）
func (*messageService) SendToChatRoom(ctx context.Context, sender model.Sender, req pb.SendMessageReq) error {
	userIds, err := cache.LargeGroupUserCache.Members(sender.AppId, req.ReceiverId)
	if err != nil {
		return err
	}

	isMember, err := cache.LargeGroupUserCache.IsMember(sender.AppId, req.ReceiverId, sender.SenderId)
	if err != nil {
		return err
	}

	if sender.SenderType == pb.SenderType_ST_USER && !isMember {
		logger.Logger.Error("not int group", zap.Int64("app_id", sender.AppId), zap.Int64("group_id", req.ReceiverId),
			zap.Int64("user_id", sender.AppId))
		return gerrors.ErrNotInGroup
	}

	var seq int64 = 0
	if req.IsPersist {
		seq, err = SeqService.GetGroupNext(ctx, sender.AppId, req.ReceiverId)
		if err != nil {
			return err
		}
		messageType, messageContent := model.PBToMessageBody(req.MessageBody)
		message := model.Message{
			AppId:          sender.AppId,
			ObjectType:     model.MessageObjectTypeGroup,
			ObjectId:       req.ReceiverId,
			RequestId:      grpclib.GetCtxRequstId(ctx),
			SenderType:     int32(sender.SenderType),
			SenderId:       sender.SenderId,
			SenderDeviceId: sender.DeviceId,
			ReceiverType:   int32(req.ReceiverType),
			ReceiverId:     req.ReceiverId,
			ToUserIds:      model.FormatUserIds(req.ToUserIds),
			Type:           messageType,
			Content:        messageContent,
			Seq:            seq,
			SendTime:       util.UnunixMilliTime(req.SendTime),
			Status:         int32(pb.MessageStatus_MS_NORMAL),
		}
		err = MessageService.Add(ctx, message)
		if err != nil {
			return err
		}
	}

	// 将消息发送给群组用户，使用读扩散
	req.IsPersist = false
	for i := range userIds {
		err = MessageService.SendToUser(ctx, sender, userIds[i].UserId, seq, req)
		if err != nil {
			return err
		}
	}
	return nil
}

// StoreAndSendToUser 将消息持久化到数据库,并且消息发送至用户
func (*messageService) SendToUser(ctx context.Context, sender model.Sender, toUserId int64, roomSeq int64, req pb.SendMessageReq) error {
	logger.Logger.Debug("message_store_send_to_user",
		zap.String("message_id", req.MessageId),
		zap.Int64("app_id", sender.AppId),
		zap.Int64("to_user_id", toUserId))

	var (
		seq = roomSeq
		err error
	)
	if req.IsPersist {
		seq, err = SeqService.GetUserNext(ctx, sender.AppId, toUserId)
		if err != nil {
			return err
		}
		messageType, messageContent := model.PBToMessageBody(req.MessageBody)
		selfMessage := model.Message{
			AppId:          sender.AppId,
			ObjectType:     model.MessageObjectTypeUser,
			ObjectId:       toUserId,
			RequestId:      grpclib.GetCtxRequstId(ctx),
			SenderType:     int32(sender.SenderType),
			SenderId:       sender.SenderId,
			SenderDeviceId: sender.DeviceId,
			ReceiverType:   int32(req.ReceiverType),
			ReceiverId:     req.ReceiverId,
			ToUserIds:      model.FormatUserIds(req.ToUserIds),
			Type:           messageType,
			Content:        messageContent,
			Seq:            seq,
			SendTime:       util.UnunixMilliTime(req.SendTime),
			Status:         int32(pb.MessageStatus_MS_NORMAL),
		}
		err = MessageService.Add(ctx, selfMessage)
		if err != nil {
			return err
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
		MessageBody:    req.MessageBody,
		Seq:            seq,
		SendTime:       req.SendTime,
		Status:         pb.MessageStatus_MS_NORMAL,
	}

	// 查询用户在线设备
	devices, err := DeviceService.ListOnlineByUserId(ctx, sender.AppId, toUserId)
	if err != nil {
		return err
	}

	for i := range devices {
		// 消息不需要投递给发送消息的设备
		if sender.DeviceId == devices[i].DeviceId {
			continue
		}

		err = MessageService.SendToDevice(ctx, devices[i], messageItem)
		if err != nil {
			return err
		}
	}
	return nil
}

// SendToDevice 将消息发送给设备
func (*messageService) SendToDevice(ctx context.Context, device model.Device, msgItem pb.MessageItem) error {
	if device.Status == model.DeviceOnLine {
		message := pb.Message{Message: &msgItem}
		_, err := rpc_cli.ConnectIntClient.DeliverMessage(grpclib.ContextWithAddr(ctx, device.ConnAddr), &pb.DeliverMessageReq{
			DeviceId: device.DeviceId, Message: &message})
		if err != nil {
			return err
		}
	}

	// todo 其他推送厂商
	return nil
}

package service

import (
	"context"
	"gim/internal/logic/domain/message/model"
	"gim/internal/logic/domain/message/repo"
	"gim/internal/logic/proxy"
	"gim/pkg/grpclib"
	"gim/pkg/grpclib/picker"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc"
	"gim/pkg/util"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const MessageLimit = 50 // 最大消息同步数量

const MaxSyncBufLen = 65536 // 最大字节数组长度

type messageService struct{}

var MessageService = new(messageService)

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
	usersResp, err := rpc.GetBusinessIntClient().GetUsers(ctx, &pb.GetUsersReq{UserIds: userIds})
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
	return repo.MessageRepo.ListBySeq(userId, seq, MessageLimit)
}

// SendToUser 将消息发送给用户
func (*messageService) SendToUser(ctx context.Context, sender *pb.Sender, toUserId int64, req *pb.SendMessageReq) (int64, error) {
	logger.Logger.Debug("SendToUser",
		zap.Int64("request_id", grpclib.GetCtxRequestId(ctx)),
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

		selfMessage := model.Message{
			UserId:       toUserId,
			RequestId:    grpclib.GetCtxRequestId(ctx),
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
		err = repo.MessageRepo.Save(selfMessage)
		if err != nil {
			logger.Sugar.Error(err)
			return 0, err
		}

		if sender.SenderType == pb.SenderType_ST_USER && sender.SenderId == toUserId {
			// 用户需要增加自己的已经同步的序列号
			err = repo.DeviceACKRepo.Set(sender.SenderId, sender.DeviceId, seq)
			if err != nil {
				return 0, err
			}
		}
	}

	message := pb.Message{
		Sender:         sender,
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
	devices, err := proxy.DeviceProxy.ListOnlineByUserId(ctx, toUserId)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}

	for i := range devices {
		// 消息不需要投递给发送消息的设备
		if sender.DeviceId == devices[i].DeviceId {
			continue
		}

		err = MessageService.SendToDevice(ctx, devices[i], &message)
		if err != nil {
			logger.Sugar.Error(err, zap.Any("SendToUser error", devices[i]), zap.Error(err))
		}
	}

	return seq, nil
}

// SendToDevice 将消息发送给设备
func (*messageService) SendToDevice(ctx context.Context, device *pb.Device, message *pb.Message) error {
	messageSend := pb.MessageSend{Message: message}
	_, err := rpc.GetConnectIntClient().DeliverMessage(picker.ContextWithAddr(ctx, device.ConnAddr), &pb.DeliverMessageReq{
		DeviceId:    device.DeviceId,
		MessageSend: &messageSend,
	})
	if err != nil {
		logger.Logger.Error("SendToDevice error", zap.Error(err))
		return err
	}

	// todo 其他推送厂商
	return nil
}

func (*messageService) AddSenderInfo(sender *pb.Sender) {
	if sender.SenderType == pb.SenderType_ST_USER {
		user, err := rpc.GetBusinessIntClient().GetUser(context.TODO(), &pb.GetUserReq{UserId: sender.SenderId})
		if err == nil && user != nil {
			sender.AvatarUrl = user.User.AvatarUrl
			sender.Nickname = user.User.Nickname
			sender.Extra = user.User.Extra
		}
	}
}

package app

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/protobuf/proto"

	deviceapp "gim/internal/logic/device/app"
	devicedomain "gim/internal/logic/device/domain"
	"gim/internal/logic/message/domain"
	"gim/internal/logic/message/repo"
	"gim/pkg/md"
	"gim/pkg/mq"
	"gim/pkg/protocol/pb/connectpb"
	pb "gim/pkg/protocol/pb/logicpb"
	"gim/pkg/rpc"
)

const pageSize = 50 // 最大消息同步数量

var MessageApp = new(messageApp)

type messageApp struct{}

// PushToUsers 发送消息
func (a *messageApp) PushToUsers(ctx context.Context, userIDs []uint64, message *connectpb.Message, isPersist bool) (uint64, error) {
	message.CreatedAt = time.Now().Unix()
	slog.Debug("SendToUser", "request_id", md.GetRequestID(ctx), "to_user_ids", userIDs)

	var messageID uint64
	if isPersist {
		msg := domain.Message{
			RequestID: md.GetRequestID(ctx),
			Command:   message.Command,
			Content:   message.Content,
			CreatedAt: time.Unix(message.CreatedAt, 0),
		}
		err := repo.MessageRepo.Save(ctx, &msg)
		if err != nil {
			return 0, err
		}
		messageID = msg.ID
	}

	for _, userID := range userIDs {
		msg := &connectpb.Message{
			RequestId: message.RequestId,
			Command:   message.Command,
			Content:   message.Content,
			CreatedAt: message.CreatedAt,
		}
		err := a.PushToUser(ctx, userID, messageID, msg, isPersist)
		if err != nil {
			slog.Error("PushToUser", "error", err, "userID", userID)
		}
	}
	return messageID, nil
}

func (a *messageApp) PushToUser(ctx context.Context, userID, messageID uint64, message *connectpb.Message, isPersist bool) error {
	slog.Debug("PushToUser", "userID", userID, "messageID", messageID, "message", message)

	if isPersist {
		userMessage := domain.UserMessage{
			UserID:    userID,
			MessageID: messageID,
		}

		err := repo.UserMessageRepo.Create(ctx, &userMessage)
		if err != nil {
			return err
		}
		message.Seq = userMessage.Seq
	}

	devices, err := deviceapp.DeviceApp.ListByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for i := range devices {
		err = a.PushToDevice(ctx, &devices[i], message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *messageApp) PushToDevice(ctx context.Context, device *devicedomain.Device, message *connectpb.Message) error {
	slog.Debug("PushToDevice", "device", device, "message", message)
	if device.IsOnline {
		request := &connectpb.PushToDevicesRequest{
			DeviceMessageList: []*connectpb.DeviceMessage{
				{
					DeviceId: 0,
					Message:  message,
				},
			},
		}

		_, err := rpc.GetConnectIntClient(device.ConnectAddr).PushToDevices(ctx, request)
		return err
	}

	// 离线推送
	return nil
}

// PushToAll 全服推送
func (*messageApp) PushToAll(ctx context.Context, req *pb.PushToAllRequest) error {
	msg := connectpb.PushAllMessage{
		Message: &connectpb.Message{
			Command:   req.Command,
			Content:   req.Content,
			CreatedAt: time.Now().Unix(),
		},
	}
	bytes, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}
	return mq.Publish(ctx, mq.PushAllTopic, bytes)
}

// Sync 消息同步
func (a *messageApp) Sync(ctx context.Context, userId, seq uint64) (*pb.SyncReply, error) {
	messages, hasMore, err := a.listByUserIdAndSeq(ctx, userId, seq)
	if err != nil {
		return nil, err
	}
	pbMessages := domain.MessagesToPB(messages)

	reply := &pb.SyncReply{Messages: pbMessages, HasMore: hasMore}
	return reply, nil
}

// listByUserIdAndSeq 查询消息
func (a *messageApp) listByUserIdAndSeq(ctx context.Context, userId, seq uint64) ([]domain.UserMessage, bool, error) {
	var err error
	if seq == 0 {
		seq, err = DeviceACKApp.getMaxByUserId(ctx, userId)
		if err != nil {
			return nil, false, err
		}
	}
	return repo.UserMessageRepo.ListBySeq(ctx, userId, seq, pageSize)
}

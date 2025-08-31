package message

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/protobuf/proto"

	"gim/internal/logic/device"
	"gim/internal/logic/message/domain"
	"gim/internal/logic/message/repo"
	"gim/pkg/md"
	"gim/pkg/mq"
	"gim/pkg/protocol/pb/connectpb"
	pb "gim/pkg/protocol/pb/logicpb"
	"gim/pkg/rpc"
)

const pageSize = 50 // 最大消息同步数量

var App = new(app)

type app struct{}

func (a *app) PushToUsersWithAny(ctx context.Context, userIDs []uint64, command connectpb.Command, any proto.Message, isPersist bool) (uint64, error) {
	bytes, err := proto.Marshal(any)
	if err != nil {
		slog.Error("PushToUser", "error", err)
		return 0, err
	}
	return a.PushToUsersWithCommand(ctx, userIDs, command, bytes, isPersist)
}

func (a *app) PushToUsersWithCommand(ctx context.Context, toUserIDs []uint64, command connectpb.Command, content []byte,
	isPersist bool) (uint64, error) {
	message := connectpb.Message{
		Command: command,
		Content: content,
		Seq:     0,
	}
	messageID, err := a.PushToUsers(ctx, toUserIDs, &message, isPersist)
	if err != nil {
		slog.Error("PushToUser", "error", err)
		return 0, err
	}
	return messageID, nil
}

// PushToUsers 发送消息
func (a *app) PushToUsers(ctx context.Context, userIDs []uint64, message *connectpb.Message, isPersist bool) (uint64, error) {
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
		err := repo.MessageRepo.Save(&msg)
		if err != nil {
			return 0, err
		}
		messageID = msg.ID
	}

	for _, userID := range userIDs {
		err := a.PushToUser(ctx, userID, messageID, *message, isPersist)
		if err != nil {
			slog.Error("PushToUser", "error", err, "userID", userID)
		}
	}
	return messageID, nil
}

func (a *app) PushToUser(ctx context.Context, userID, messageID uint64, message connectpb.Message, isPersist bool) error {
	slog.Debug("PushToUser", "userID", userID, "messageID", messageID, "message", message)
	var (
		seq uint64
		err error
	)
	if isPersist {
		seq, err = repo.SeqRepo.Incr(repo.SeqObjectTypeUser, userID)
		if err != nil {
			return err
		}

		userMessage := domain.UserMessage{
			UserID:    userID,
			Seq:       seq,
			MessageID: messageID,
		}

		err = repo.UserMessageRepo.Create(&userMessage)
		if err != nil {
			return err
		}
	}

	message.Seq = seq
	devices, err := device.App.ListByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for i := range devices {
		err = a.PushToDevice(ctx, &devices[i], &message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *app) PushToDevice(ctx context.Context, device *device.Device, message *connectpb.Message) error {
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
func (*app) PushToAll(ctx context.Context, req *pb.PushToAllRequest) error {
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
	return mq.Publish(mq.PushAllTopic, bytes)
}

// Sync 消息同步
func (a *app) Sync(ctx context.Context, userId, seq uint64) (*connectpb.SyncReply, error) {
	messages, hasMore, err := a.listByUserIdAndSeq(ctx, userId, seq)
	if err != nil {
		return nil, err
	}
	pbMessages := domain.MessagesToPB(messages)

	reply := &connectpb.SyncReply{Messages: pbMessages, HasMore: hasMore}
	return reply, nil
}

// listByUserIdAndSeq 查询消息
func (a *app) listByUserIdAndSeq(ctx context.Context, userId, seq uint64) ([]domain.UserMessage, bool, error) {
	var err error
	if seq == 0 {
		seq, err = a.getMaxByUserId(ctx, userId)
		if err != nil {
			return nil, false, err
		}
	}
	return repo.UserMessageRepo.ListBySeq(userId, seq, pageSize)
}

// getMaxByUserId 根据用户id获取最大ack
func (*app) getMaxByUserId(ctx context.Context, userId uint64) (uint64, error) {
	acks, err := repo.DeviceACKRepo.Get(userId)
	if err != nil {
		return 0, err
	}

	var max uint64 = 0
	for i := range acks {
		if acks[i] > max {
			max = acks[i]
		}
	}
	return max, nil
}

// MessageAck 收到消息回执
func (*app) MessageAck(ctx context.Context, userId, deviceId, ack uint64) error {
	if ack <= 0 {
		return nil
	}
	return repo.DeviceACKRepo.Set(userId, deviceId, ack)
}

package app

import (
	"context"
	"log/slog"
	"time"

	"gim/config"
	deviceapp "gim/internal/logic/device/app"
	devicedomain "gim/internal/logic/device/domain"
	"gim/internal/logic/message/domain"
	"gim/internal/logic/message/repo"
	"gim/pkg/md"
	"gim/pkg/protocol/pb/connectpb"
	"gim/pkg/rpc"
)

type dispatcher struct {
	ctx       context.Context
	message   *connectpb.Message
	isPersist bool
	userIDs   []uint64

	messageID      uint64
	userMessages   []userMessage
	deviceMessages map[string][]*connectpb.DeviceMessage
}

type userMessage struct {
	*domain.UserMessage
	*connectpb.Message
}

func newDispatcher(ctx context.Context, message *connectpb.Message, isPersist bool, userIDs []uint64) *dispatcher {
	return &dispatcher{
		ctx:            ctx,
		message:        message,
		isPersist:      isPersist,
		userIDs:        userIDs,
		deviceMessages: make(map[string][]*connectpb.DeviceMessage),
	}
}

func (d *dispatcher) dispatch() error {
	if d.isPersist {
		msg := domain.Message{
			RequestID: md.GetRequestID(d.ctx),
			Command:   d.message.Command,
			Content:   d.message.Content,
			CreatedAt: time.Unix(d.message.CreatedAt, 0),
		}
		err := repo.MessageRepo.Save(d.ctx, &msg)
		if err != nil {
			return err
		}
		d.messageID = msg.ID
	}

	for _, userID := range d.userIDs {
		var seq uint64
		var err error
		if d.isPersist {
			seq, err = repo.SeqRepo.Incr(d.ctx, userID)
			if err != nil {
				slog.ErrorContext(d.ctx, "messageDispatcher dispatch SeqRepo Incr", "error", err, "user_id", userID)
				continue
			}
		}

		userMessage := userMessage{
			UserMessage: &domain.UserMessage{
				UserID:    userID,
				Seq:       seq,
				MessageID: d.messageID,
			},
			Message: &connectpb.Message{
				Command:   d.message.Command,
				Content:   d.message.Content,
				CreatedAt: d.message.CreatedAt,
				Seq:       seq,
			},
		}
		d.userMessages = append(d.userMessages, userMessage)
	}

	for _, userMessage := range d.userMessages {
		devices, err := deviceapp.DeviceApp.ListByUserID(d.ctx, userMessage.UserID)
		if err != nil {
			slog.ErrorContext(d.ctx, "messageDispatcher dispatch ListByUserID", "error", err, "user_id", userMessage.UserID)
			continue
		}

		for _, device := range devices {
			if device.Status == devicedomain.StatusOnline {
				d.appendDeviceMessage(device.ConnectIP, device.ID, userMessage.Message)
			} else {
				// 离线推送
			}
		}
	}
	return nil
}

func (d *dispatcher) appendDeviceMessage(ip string, deviceID uint64, message *connectpb.Message) {
	d.deviceMessages[ip] = append(d.deviceMessages[ip], &connectpb.DeviceMessage{
		DeviceId: deviceID,
		Message:  message,
	})
}

func (d *dispatcher) flush() error {
	if d.isPersist {
		userMessages := make([]*domain.UserMessage, 0, len(d.userMessages))

		for _, userMessage := range d.userMessages {
			userMessages = append(userMessages, userMessage.UserMessage)
		}

		err := repo.UserMessageRepo.BatchCreate(d.ctx, userMessages)
		if err != nil {
			return err
		}
	}

	for ip, dms := range d.deviceMessages {
		addr := ip + config.GrpcListenAddr
		pushToDevices(d.ctx, addr, dms)
	}
	return nil
}

func pushToDevices(ctx context.Context, addr string, deviceMessages []*connectpb.DeviceMessage) {
	for i := 0; i < len(deviceMessages); i += 100 {
		end := min(i+100, len(deviceMessages))
		_, err := rpc.GetConnectIntClient(addr).PushToDevices(ctx, &connectpb.PushToDevicesRequest{
			DeviceMessages: deviceMessages[i:end],
		})
		if err != nil {
			slog.Error("pushToDevices", "error", err, "addr", addr)
		}
	}
}

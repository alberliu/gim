package app

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/protobuf/proto"

	"gim/internal/logic/message/domain"
	"gim/internal/logic/message/repo"
	"gim/pkg/md"
	"gim/pkg/mq"
	"gim/pkg/protocol/pb/connectpb"
	pb "gim/pkg/protocol/pb/logicpb"
)

const pageSize = 50 // 最大消息同步数量

var MessageApp = new(messageApp)

type messageApp struct{}

// PushToUsers 发送消息
func (a *messageApp) PushToUsers(ctx context.Context, userIDs []uint64, message *connectpb.Message, isPersist bool) (uint64, error) {
	message.CreatedAt = time.Now().Unix()
	slog.DebugContext(ctx, "PushToUsers", "request_id", md.GetRequestID(ctx), "to_user_ids", userIDs)

	dispatcher := newDispatcher(ctx, message, isPersist, userIDs)
	err := dispatcher.dispatch()
	if err != nil {
		return 0, err
	}
	err = dispatcher.flush()
	if err != nil {
		return 0, err
	}
	return dispatcher.messageID, nil
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

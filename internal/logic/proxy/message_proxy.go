package proxy

import (
	"context"
	"gim/pkg/logger"
	"gim/pkg/protocol/pb"
	"gim/pkg/util"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var MessageProxy messageProxy

type messageProxy interface {
	SendToUser(ctx context.Context, fromDeviceID, toUserID int64, message *pb.Message, isPersist bool) (int64, error)
}

func PushToUserBytes(ctx context.Context, toUserID int64, code int32, bytes []byte, isPersist bool) (int64, error) {
	message := pb.Message{
		Code:     code,
		Content:  bytes,
		SendTime: util.UnixMilliTime(time.Now()),
	}
	seq, err := MessageProxy.SendToUser(ctx, 0, toUserID, &message, isPersist)
	if err != nil {
		logger.Logger.Error("PushToUser", zap.Error(err))
		return 0, err
	}
	return seq, nil
}

func PushToUser(ctx context.Context, toUserID int64, code pb.PushCode, msg proto.Message, isPersist bool) (int64, error) {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		logger.Logger.Error("PushToUser", zap.Error(err))
		return 0, err
	}
	return PushToUserBytes(ctx, toUserID, int32(code), bytes, isPersist)
}

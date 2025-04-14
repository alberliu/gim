package service

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	"gim/pkg/gerrors"
	"gim/pkg/mq"
	"gim/pkg/protocol/pb"
	"gim/pkg/util"
)

type pushService struct{}

var PushService = new(pushService)

// PushAll 全服推送
func (s *pushService) PushAll(ctx context.Context, req *pb.PushAllReq) error {
	msg := pb.PushAllMsg{
		Message: &pb.Message{
			Code:     req.Code,
			Content:  req.Content,
			SendTime: util.UnixMilliTime(time.Now()),
		},
	}
	bytes, err := proto.Marshal(&msg)
	if err != nil {
		return gerrors.WrapError(err)
	}
	err = mq.Publish(mq.PushAllTopic, bytes)
	if err != nil {
		return err
	}
	return nil
}

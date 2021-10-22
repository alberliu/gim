package proxy

import (
	"context"
	"gim/pkg/pb"

	"google.golang.org/protobuf/proto"
)

var MessageProxy messageProxy

type messageProxy interface {
	SendToUser(ctx context.Context, sender *pb.Sender, toUserId int64, req *pb.SendMessageReq) (int64, error)
	PushToUser(ctx context.Context, userId int64, code pb.PushCode, message proto.Message, isPersist bool) error
}

package api

import (
	"context"

	"gim/internal/logic/message/app"
	"gim/pkg/md"
	pb "gim/pkg/protocol/pb/logicpb"
)

type MessageExtService struct {
	pb.UnsafeMessageExtServiceServer
}

func (m MessageExtService) Sync(ctx context.Context, request *pb.SyncRequest) (*pb.SyncReply, error) {
	userID := md.GetUserID(ctx)
	return app.MessageApp.Sync(ctx, userID, request.Seq)
}

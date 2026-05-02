package api

import (
	"context"

	pb "gim/pkg/gen/proto/businesspb"
	"gim/pkg/gen/proto/logicpb"
	"gim/pkg/rpc"
)

type GroupExtService struct {
	pb.UnsafeGroupExtServiceServer
}

func (*GroupExtService) Create(ctx context.Context, req *logicpb.GroupCreateRequest) (*logicpb.GroupCreateReply, error) {
	return rpc.GetGroupIntClient().Create(ctx, req)
}

func (*GroupExtService) Get(ctx context.Context, req *logicpb.GroupGetRequest) (*logicpb.GroupGetReply, error) {
	return rpc.GetGroupIntClient().Get(ctx, req)
}

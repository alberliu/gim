package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"gim/internal/business/domain/user/app"
	"gim/pkg/protocol/pb"
)

type BusinessIntServer struct {
	pb.UnsafeBusinessIntServer
}

func (*BusinessIntServer) Auth(ctx context.Context, req *pb.AuthReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, app.AuthApp.Auth(ctx, req.UserId, req.DeviceId, req.Token)
}

func (*BusinessIntServer) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserResp, error) {
	user, err := app.UserApp.Get(ctx, req.UserId)
	return &pb.GetUserResp{User: user}, err
}

func (*BusinessIntServer) GetUsers(ctx context.Context, req *pb.GetUsersReq) (*pb.GetUsersResp, error) {
	var userIds = make([]int64, 0, len(req.UserIds))
	for k := range req.UserIds {
		userIds = append(userIds, k)
	}

	users, err := app.UserApp.GetByIds(ctx, userIds)
	return &pb.GetUsersResp{Users: users}, err
}

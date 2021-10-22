package api

import (
	"context"
	"gim/internal/business/app"
	"gim/pkg/pb"
)

type BusinessIntServer struct{}

func (*BusinessIntServer) Auth(ctx context.Context, req *pb.AuthReq) (*pb.Empty, error) {
	return &pb.Empty{}, app.AuthApp.Auth(ctx, req.UserId, req.DeviceId, req.Token)
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

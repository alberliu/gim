package api

import (
	"context"
	"gim/internal/business/service"
	"gim/pkg/pb"
)

type BusinessIntServer struct{}

func (*BusinessIntServer) Auth(ctx context.Context, req *pb.AuthReq) (*pb.Empty, error) {
	return &pb.Empty{}, service.AuthService.Auth(ctx, req.UserId, req.DeviceId, req.Token)
}

func (*BusinessIntServer) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserResp, error) {
	user, err := service.UserService.Get(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResp{User: user.ToProto()}, nil
}

func (*BusinessIntServer) GetUsers(ctx context.Context, req *pb.GetUsersReq) (*pb.GetUsersResp, error) {
	var userIds = make([]int64, 0, len(req.UserIds))
	for k := range req.UserIds {
		userIds = append(userIds, k)
	}

	users, err := service.UserService.GetByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	pbUsers := make(map[int64]*pb.User, len(users))
	for i := range users {
		pbUsers[users[i].Id] = users[i].ToProto()
	}

	return &pb.GetUsersResp{Users: pbUsers}, nil
}

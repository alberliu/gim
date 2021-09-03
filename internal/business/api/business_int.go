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
	return &pb.GetUserResp{User: &pb.User{
		UserId:     user.Id,
		Nickname:   user.Nickname,
		Sex:        user.Sex,
		AvatarUrl:  user.AvatarUrl,
		Extra:      user.Extra,
		CreateTime: user.CreateTime.Unix(),
		UpdateTime: user.UpdateTime.Unix(),
	}}, nil
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
		pbUsers[users[i].Id] = &pb.User{
			UserId:     users[i].Id,
			Nickname:   users[i].Nickname,
			Sex:        users[i].Sex,
			AvatarUrl:  users[i].AvatarUrl,
			Extra:      users[i].Extra,
			CreateTime: users[i].CreateTime.Unix(),
			UpdateTime: users[i].UpdateTime.Unix(),
		}
	}

	return &pb.GetUsersResp{Users: pbUsers}, nil
}

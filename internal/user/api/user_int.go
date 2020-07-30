package api

import (
	"context"
	"gim/internal/user/service"
	"gim/pkg/pb"
)

type UserIntServer struct{}

func (*UserIntServer) Auth(ctx context.Context, req *pb.AuthReq) (*pb.AuthResp, error) {
	return &pb.AuthResp{}, service.AuthService.Auth(ctx, req.UserId, req.DeviceId, req.Token)
}

func (*UserIntServer) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserResp, error) {
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

func (*UserIntServer) GetUsers(ctx context.Context, req *pb.GetUsersReq) (*pb.GetUsersResp, error) {
	users, err := service.UserService.GetByIds(ctx, req.UserIds)
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

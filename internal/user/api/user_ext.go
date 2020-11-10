package api

import (
	"context"
	"gim/internal/user/dao"
	"gim/internal/user/model"
	"gim/internal/user/service"
	"gim/pkg/grpclib"
	"gim/pkg/pb"
)

type UserExtServer struct{}

func (s *UserExtServer) SignIn(ctx context.Context, req *pb.SignInReq) (*pb.SignInResp, error) {
	userId, token, err := service.AuthService.SignIn(ctx, req.PhoneNumber, req.Code, req.DeviceId)
	if err != nil {
		return nil, err
	}
	return &pb.SignInResp{
		UserId: userId,
		Token:  token,
	}, nil
}

func (s *UserExtServer) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}
	user, err := service.UserService.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResp{
		User: user.ToProto(),
	}, nil
}

func (s *UserExtServer) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.UpdateUserResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserResp{}, service.UserService.Update(ctx, model.User{
		Id:        userId,
		Nickname:  req.Nickname,
		Sex:       req.Sex,
		AvatarUrl: req.AvatarUrl,
		Extra:     req.Extra,
	})
}

func (s *UserExtServer) SearchUser(ctx context.Context, req *pb.SearchUserReq) (*pb.SearchUserResp, error) {
	users, err := dao.UserDao.Search(req.Key)
	if err != nil {
		return nil, err
	}
	pbUsers := make([]*pb.User, 0, len(users))
	for i := range users {
		pbUsers = append(pbUsers, users[i].ToProto())
	}
	return &pb.SearchUserResp{Users: pbUsers}, nil
}

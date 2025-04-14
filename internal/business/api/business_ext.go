package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"gim/internal/business/domain/user/app"
	"gim/pkg/grpclib"
	"gim/pkg/protocol/pb"
)

type BusinessExtServer struct {
	pb.UnsafeBusinessExtServer
}

func (s *BusinessExtServer) SignIn(ctx context.Context, req *pb.SignInReq) (*pb.SignInResp, error) {
	isNew, userId, token, err := app.AuthApp.SignIn(ctx, req.PhoneNumber, req.Code, req.DeviceId)
	if err != nil {
		return nil, err
	}
	return &pb.SignInResp{
		IsNew:  isNew,
		UserId: userId,
		Token:  token,
	}, nil
}

func (s *BusinessExtServer) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	user, err := app.UserApp.Get(ctx, userId)
	return &pb.GetUserResp{User: user}, err
}

func (s *BusinessExtServer) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*emptypb.Empty, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	return new(emptypb.Empty), app.UserApp.Update(ctx, userId, req)
}

func (s *BusinessExtServer) SearchUser(ctx context.Context, req *pb.SearchUserReq) (*pb.SearchUserResp, error) {
	users, err := app.UserApp.Search(ctx, req.Key)
	return &pb.SearchUserResp{Users: users}, err
}

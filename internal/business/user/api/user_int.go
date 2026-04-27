package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"gim/internal/business/user/app"
	pb "gim/pkg/protocol/pb/businesspb"
)

type UserIntService struct {
	pb.UnsafeUserIntServiceServer
}

func (*UserIntService) Auth(ctx context.Context, req *pb.AuthRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, app.AuthApp.Auth(ctx, req.UserId, req.DeviceId, req.Token)
}

func (*UserIntService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	user, err := app.UserApp.Get(ctx, req.UserId)
	return &pb.GetUserReply{User: user}, err
}

func (*UserIntService) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersReply, error) {
	users, err := app.UserApp.GetUsers(ctx, req.UserIds)
	return &pb.GetUsersReply{Users: users}, err
}

package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"gim/internal/business/user/app"
	"gim/pkg/md"
	pb "gim/pkg/protocol/pb/businesspb"
)

type UserExtService struct {
	pb.UnsafeUserExtServiceServer
}

func (s *UserExtService) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInReply, error) {
	return app.AuthApp.SignIn(ctx, req)
}

func (s *UserExtService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	userID := md.GetUserID(ctx)

	user, err := app.UserApp.Get(ctx, userID)
	return &pb.GetUserReply{User: user}, err
}

func (s *UserExtService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*emptypb.Empty, error) {
	userID := md.GetUserID(ctx)

	return new(emptypb.Empty), app.UserApp.Update(ctx, userID, req)
}

func (s *UserExtService) SearchUser(ctx context.Context, req *pb.SearchUserRequest) (*pb.SearchUserReply, error) {
	users, err := app.UserApp.Search(ctx, req.Key)
	return &pb.SearchUserReply{Users: users}, err
}

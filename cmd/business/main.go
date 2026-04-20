package main

import (
	"google.golang.org/grpc"

	friendapi "gim/internal/business/friend/api"
	userapi "gim/internal/business/user/api"
	"gim/pkg/logger"
	pb "gim/pkg/protocol/pb/businesspb"
	"gim/pkg/server"
)

func main() {
	logger.Init("user")

	server.RunGRPCServer(func(server *grpc.Server) {
		pb.RegisterUserIntServiceServer(server, &userapi.UserIntService{})
		pb.RegisterUserExtServiceServer(server, &userapi.UserExtService{})
		pb.RegisterFriendExtServiceServer(server, &friendapi.FriendExtService{})
	})

	server.WaitForShutdown()
}

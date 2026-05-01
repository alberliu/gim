package main

import (
	"google.golang.org/grpc"

	friendapi "gim/internal/business/friend/api"
	messageapi "gim/internal/business/message/api"
	userapi "gim/internal/business/user/api"
	"gim/pkg/db"
	pb "gim/pkg/gen/proto/businesspb"
	"gim/pkg/logger"
	"gim/pkg/server"
)

func main() {
	logger.Init()
	db.Init()

	server.RunGRPCServer(func(server *grpc.Server) {
		pb.RegisterUserIntServiceServer(server, &userapi.UserIntService{})
		pb.RegisterUserExtServiceServer(server, &userapi.UserExtService{})
		pb.RegisterFriendExtServiceServer(server, &friendapi.FriendExtService{})
		pb.RegisterMessageExtServiceServer(server, &messageapi.MessageExtService{})
	})

	server.WaitForShutdown()
}

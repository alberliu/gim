package main

import (
	"google.golang.org/grpc"

	deviceapi "gim/internal/logic/device/api"
	groupapi "gim/internal/logic/group/api"
	messageapi "gim/internal/logic/message/api"
	"gim/internal/logic/room"
	"gim/pkg/logger"
	pb "gim/pkg/protocol/pb/logicpb"
	"gim/pkg/server"
)

func main() {
	logger.Init("logic")

	server.RunGRPCServer(func(server *grpc.Server) {
		pb.RegisterDeviceIntServiceServer(server, &deviceapi.DeviceIntService{})
		pb.RegisterMessageExtServiceServer(server, &messageapi.MessageExtService{})
		pb.RegisterMessageIntServiceServer(server, &messageapi.MessageIntService{})
		pb.RegisterGroupIntServiceServer(server, &groupapi.GroupIntService{})
		pb.RegisterRoomIntServiceServer(server, &room.RoomIntService{})

	})

	server.WaitForShutdown()
}

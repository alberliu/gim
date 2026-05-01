package main

import (
	"google.golang.org/grpc"

	deviceapi "gim/internal/logic/device/api"
	groupapi "gim/internal/logic/group/api"
	messageapi "gim/internal/logic/message/api"
	"gim/internal/logic/room"
	"gim/pkg/db"
	pb "gim/pkg/gen/proto/logicpb"
	"gim/pkg/logger"
	"gim/pkg/server"
)

func main() {
	logger.Init()
	db.Init()

	server.RunGRPCServer(func(server *grpc.Server) {
		pb.RegisterDeviceIntServiceServer(server, &deviceapi.DeviceIntService{})
		pb.RegisterMessageExtServiceServer(server, &messageapi.MessageExtService{})
		pb.RegisterMessageIntServiceServer(server, &messageapi.MessageIntService{})
		pb.RegisterGroupIntServiceServer(server, &groupapi.GroupIntService{})
		pb.RegisterRoomIntServiceServer(server, &room.RoomIntService{})
	})

	server.WaitForShutdown()
}

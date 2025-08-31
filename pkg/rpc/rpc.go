package rpc

import (
	"context"
	"sync"

	"gim/config"
	"gim/pkg/protocol/pb/businesspb"
	"gim/pkg/protocol/pb/connectpb"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/ugrpc"
)

var connectIntClients sync.Map

var (
	deviceIntClient  logicpb.DeviceIntServiceClient
	messageIntClient logicpb.MessageIntServiceClient
	roomIntClient    logicpb.RoomIntServiceClient
	userIntClient    businesspb.UserIntServiceClient
)

func GetConnectIntClient(addr string) connectpb.ConnectIntServiceClient {
	value, ok := connectIntClients.Load(addr)
	if ok {
		return value.(connectpb.ConnectIntServiceClient)
	}

	conn := ugrpc.NewClient(addr)
	client := connectpb.NewConnectIntServiceClient(conn)
	connectIntClients.Store(addr, client)
	return client
}

func GetDeviceIntClient() logicpb.DeviceIntServiceClient {
	if deviceIntClient == nil {
		deviceIntClient = config.Config.DeviceIntClientBuilder()
	}
	return deviceIntClient
}

func GetMessageIntClient() logicpb.MessageIntServiceClient {
	if messageIntClient == nil {
		messageIntClient = config.Config.MessageIntClientBuilder()
	}
	return messageIntClient
}

func GetRoomIntClient() logicpb.RoomIntServiceClient {
	if roomIntClient == nil {
		roomIntClient = config.Config.RoomIntClientBuilder()
	}
	return roomIntClient
}

func GetUserIntClient() businesspb.UserIntServiceClient {
	if userIntClient == nil {
		userIntClient = config.Config.UserIntClientBuilder()
	}
	return userIntClient
}

func GetUser(deviceID, userID uint64) (*logicpb.User, error) {
	user, err := GetUserIntClient().GetUser(context.TODO(), &businesspb.GetUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return &logicpb.User{
		UserId:    userID,
		DeviceId:  deviceID,
		AvatarUrl: user.User.AvatarUrl,
		Nickname:  user.User.Nickname,
		Extra:     user.User.Extra,
	}, nil
}

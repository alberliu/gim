package rpc

import (
	"context"

	"gim/config"
	"gim/pkg/protocol/pb/connectpb"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/protocol/pb/userpb"
)

var (
	connectIntClient connectpb.ConnectIntServiceClient
	deviceIntClient  logicpb.DeviceIntServiceClient
	messageIntClient logicpb.MessageIntServiceClient
	roomIntClient    logicpb.RoomIntServiceClient
	userIntClient    userpb.UserIntServiceClient
)

func GetConnectIntClient() connectpb.ConnectIntServiceClient {
	if connectIntClient == nil {
		connectIntClient = config.Config.ConnectIntClientBuilder()
	}
	return connectIntClient
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

func GetUserIntClient() userpb.UserIntServiceClient {
	if userIntClient == nil {
		userIntClient = config.Config.UserIntClientBuilder()
	}
	return userIntClient
}

func GetSender(deviceID, userID uint64) (*logicpb.Sender, error) {
	user, err := GetUserIntClient().GetUser(context.TODO(), &userpb.GetUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return &logicpb.Sender{
		UserId:    userID,
		DeviceId:  deviceID,
		AvatarUrl: user.User.AvatarUrl,
		Nickname:  user.User.Nickname,
		Extra:     user.User.Extra,
	}, nil
}

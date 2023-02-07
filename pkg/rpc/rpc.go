package rpc

import (
	"context"
	"gim/config"
	"gim/pkg/protocol/pb"
)

var (
	connectIntClient  pb.ConnectIntClient
	logicIntClient    pb.LogicIntClient
	businessIntClient pb.BusinessIntClient
)

func GetConnectIntClient() pb.ConnectIntClient {
	if connectIntClient == nil {
		connectIntClient = config.Config.ConnectIntClientBuilder()
	}
	return connectIntClient
}

func GetLogicIntClient() pb.LogicIntClient {
	if logicIntClient == nil {
		logicIntClient = config.Config.LogicIntClientBuilder()
	}
	return logicIntClient
}

func GetBusinessIntClient() pb.BusinessIntClient {
	if businessIntClient == nil {
		businessIntClient = config.Config.BusinessIntClientBuilder()
	}
	return businessIntClient
}

func GetSender(deviceID, userID int64) (*pb.Sender, error) {
	user, err := GetBusinessIntClient().GetUser(context.TODO(), &pb.GetUserReq{UserId: userID})
	if err != nil {
		return nil, err
	}
	return &pb.Sender{
		UserId:    userID,
		DeviceId:  deviceID,
		AvatarUrl: user.User.AvatarUrl,
		Nickname:  user.User.Nickname,
		Extra:     user.User.Extra,
	}, nil
}

package local

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gim/config"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/protocol/pb/userpb"
)

func Init() {
	logicConn, err := grpc.NewClient("127.0.0.1:8010", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	config.Config.DeviceIntClientBuilder = func() logicpb.DeviceIntServiceClient {
		return logicpb.NewDeviceIntServiceClient(logicConn)
	}

	config.Config.MessageIntClientBuilder = func() logicpb.MessageIntServiceClient {
		return logicpb.NewMessageIntServiceClient(logicConn)
	}

	config.Config.RoomIntClientBuilder = func() logicpb.RoomIntServiceClient {
		return logicpb.NewRoomIntServiceClient(logicConn)
	}

	config.Config.UserIntClientBuilder = func() userpb.UserIntServiceClient {
		conn, err := grpc.NewClient("127.0.0.1:8020", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		return userpb.NewUserIntServiceClient(conn)
	}
}

package rpc

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/protobuf/proto"

	"gim/config"
	"gim/pkg/protocol/pb/businesspb"
	"gim/pkg/protocol/pb/connectpb"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/ugrpc"
)

func getTarget(server string) string {
	switch config.ENV {
	case config.EnvCompose:
		return fmt.Sprintf("dns:///%s:8000", server)
	case config.EnvK8s:
		return fmt.Sprintf("k8s:///%s:8000", server)
	default:
		panic("unknown env")
	}
}

var connectIntClients sync.Map

var (
	logicConn    = ugrpc.NewClient(getTarget("logic"))
	businessConn = ugrpc.NewClient(getTarget("business"))
)

var (
	deviceIntClient  = logicpb.NewDeviceIntServiceClient(logicConn)
	messageIntClient = logicpb.NewMessageIntServiceClient(logicConn)
	roomIntClient    = logicpb.NewRoomIntServiceClient(logicConn)
	userIntClient    = businesspb.NewUserIntServiceClient(businessConn)
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
	return deviceIntClient
}

func GetMessageIntClient() logicpb.MessageIntServiceClient {
	return messageIntClient
}

func GetRoomIntClient() logicpb.RoomIntServiceClient {
	return roomIntClient
}

func GetUserIntClient() businesspb.UserIntServiceClient {
	return userIntClient
}

type PushRequest struct {
	UserIDs   []uint64
	Command   connectpb.MessageCommand
	Message   proto.Message
	IsPersist bool
}

func PushToUsers(ctx context.Context, request PushRequest) (*logicpb.PushToUsersReply, error) {
	content, err := proto.Marshal(request.Message)
	if err != nil {
		return nil, err
	}
	return GetMessageIntClient().PushToUsers(ctx, &logicpb.PushToUsersRequest{
		UserIds:   request.UserIDs,
		Command:   request.Command,
		Content:   content,
		IsPersist: request.IsPersist,
	})

}

package rpc

import (
	"context"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"gim/pkg/gen/proto/businesspb"
	"gim/pkg/gen/proto/connectpb"
	"gim/pkg/gen/proto/logicpb"
	"gim/pkg/ugrpc"
)

const Timeout = time.Second * 5

var connectIntClients sync.Map

var (
	logicConn    = ugrpc.NewClient(ugrpc.GetTarget("logic"))
	businessConn = ugrpc.NewClient(ugrpc.GetTarget("business"))
)

var (
	deviceIntClient  = logicpb.NewDeviceIntServiceClient(logicConn)
	messageIntClient = logicpb.NewMessageIntServiceClient(logicConn)
	groupIntClient   = logicpb.NewGroupIntServiceClient(logicConn)
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

func GetGroupIntClient() logicpb.GroupIntServiceClient {
	return groupIntClient
}

func GetRoomIntClient() logicpb.RoomIntServiceClient {
	return roomIntClient
}

func GetUserIntClient() businesspb.UserIntServiceClient {
	return userIntClient
}

type PushRequest struct {
	FromUserID uint64
	UserIDs    []uint64
	Command    connectpb.MessageCommand
	Message    proto.Message
	IsPersist  bool
}

func PushToUsers(ctx context.Context, request PushRequest) (*logicpb.PushToUsersReply, error) {
	content, err := proto.Marshal(request.Message)
	if err != nil {
		return nil, err
	}
	return GetMessageIntClient().PushToUsers(ctx, &logicpb.PushToUsersRequest{
		FromUserId: request.FromUserID,
		UserIds:    request.UserIDs,
		Command:    request.Command,
		Content:    content,
		IsPersist:  request.IsPersist,
	})

}

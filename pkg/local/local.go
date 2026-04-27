package local

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gim/pkg/protocol/pb/businesspb"
)

const BusinessServerAddr = "127.0.0.1:8080"

func GetConn() *grpc.ClientConn {
	conn, err := grpc.NewClient(BusinessServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return conn
}

func GetUserExtServiceClient() businesspb.UserExtServiceClient {
	return businesspb.NewUserExtServiceClient(GetConn())
}

func GetMessageExtClient() businesspb.MessageExtServiceClient {
	return businesspb.NewMessageExtServiceClient(GetConn())
}

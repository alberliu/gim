package api

import (
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "gim/pkg/protocol/pb/businesspb"
)

func getUserIntClient() pb.UserIntServiceClient {
	conn, err := grpc.NewClient("127.0.0.1:8020", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return pb.NewUserIntServiceClient(conn)
}

func TestUserIntServer_Auth(t *testing.T) {
	_, err := getUserIntClient().Auth(getCtx(), &pb.AuthRequest{
		UserId: 2,
		Token:  "0",
	})
	t.Log(err)
}

func TestUserIntServer_GetUsers(t *testing.T) {
	reply, err := getUserIntClient().GetUsers(getCtx(), &pb.GetUsersRequest{
		UserIds: map[uint64]int32{1: 0, 2: 0, 3: 0},
	})
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range reply.Users {
		t.Log(k, v)
	}
}

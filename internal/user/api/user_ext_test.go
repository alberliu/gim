package api

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "gim/pkg/protocol/pb/userpb"
)

func getUserExtServiceClient() pb.UserExtServiceClient {
	conn, err := grpc.NewClient("127.0.0.1:8020", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return pb.NewUserExtServiceClient(conn)
}

func getCtx() context.Context {
	token := "0"
	return metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
		"user_id", "1",
		"device_id", "1",
		"token", token,
		"request_id", strconv.FormatInt(time.Now().UnixNano(), 10)))
}

func TestUserExtServer_SignIn(t *testing.T) {
	reply, err := getUserExtServiceClient().SignIn(context.TODO(), &pb.SignInRequest{
		PhoneNumber: "22222222222",
		Code:        "0",
		DeviceId:    2,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply)
}

func TestUserExtServer_GetUser(t *testing.T) {
	reply, err := getUserExtServiceClient().GetUser(getCtx(), &pb.GetUserRequest{UserId: 1})
	if err != nil {
		fmt.Println(err)
	}
	t.Log(reply)
}

func TestUserExtService_SearchUser(t *testing.T) {
	reply, err := getUserExtServiceClient().SearchUser(getCtx(), &pb.SearchUserRequest{Key: "1"})
	if err != nil {
		fmt.Println(err)
	}
	t.Log(reply)
}

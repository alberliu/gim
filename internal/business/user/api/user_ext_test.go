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

	pb "gim/pkg/protocol/pb/businesspb"
	"gim/pkg/protocol/pb/logicpb"
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
		PhoneNumber: "2",
		Code:        "0",
		Device: &logicpb.Device{
			Id:            2,
			Type:          logicpb.DeviceType_DT_ANDROID,
			Brand:         "xiaomi",
			Model:         "xiaomi 15",
			SystemVersion: "15.0.0",
			SdkVersion:    "1.0.0",
			BranchPushId:  "xiaomi push id",
		},
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

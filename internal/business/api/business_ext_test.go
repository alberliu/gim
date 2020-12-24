package api

import (
	"context"
	"fmt"
	"gim/pkg/pb"
	"strconv"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func getUserExtClient() pb.UserExtClient {
	conn, err := grpc.Dial("127.0.0.1:50301", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return pb.NewUserExtClient(conn)
}

func getCtx() context.Context {
	token := "0"
	return metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
		"user_id", "3",
		"device_id", "1",
		"token", token,
		"request_id", strconv.FormatInt(time.Now().UnixNano(), 10)))
}

func TestUserExtServer_SignIn(t *testing.T) {
	resp, err := getUserExtClient().SignIn(getCtx(), &pb.SignInReq{
		PhoneNumber: "18599866595",
		Code:        "1",
		DeviceId:    1,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", resp)
}

func TestUserExtServer_GetUser(t *testing.T) {
	resp, err := getUserExtClient().GetUser(getCtx(), &pb.GetUserReq{UserId: 1})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", resp)
}

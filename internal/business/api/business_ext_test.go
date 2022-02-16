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

func getBusinessExtClient() pb.BusinessExtClient {
	conn, err := grpc.Dial("111.229.238.28:50200", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return pb.NewBusinessExtClient(conn)
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
	resp, err := getBusinessExtClient().SignIn(getCtx(), &pb.SignInReq{
		PhoneNumber: "11111111111",
		Code:        "1",
		DeviceId:    1,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", resp)
}

func TestUserExtServer_GetUser(t *testing.T) {
	resp, err := getBusinessExtClient().GetUser(getCtx(), &pb.GetUserReq{UserId: 1})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", resp)
}

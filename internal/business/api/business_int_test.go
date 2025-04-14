package api

import (
	"fmt"
	"testing"

	"google.golang.org/grpc"

	"gim/pkg/protocol/pb"
)

func getBusinessIntClient() pb.BusinessIntClient {
	conn, err := grpc.Dial("localhost:50300", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return pb.NewBusinessIntClient(conn)
}

func TestUserIntServer_Auth(t *testing.T) {
	_, err := getBusinessIntClient().Auth(getCtx(), &pb.AuthReq{
		UserId:   3,
		DeviceId: 1,
		Token:    "0",
	})
	fmt.Println(err)
}

func TestUserIntServer_GetUsers(t *testing.T) {
	resp, err := getBusinessIntClient().GetUsers(getCtx(), &pb.GetUsersReq{
		UserIds: map[int64]int32{1: 0, 2: 0, 3: 0},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, v := range resp.Users {
		fmt.Printf("%+-5v  %+v\n", k, v)
	}
}

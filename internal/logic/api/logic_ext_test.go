package api

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"gim/pkg/protocol/pb"
	"gim/pkg/util"
)

func getLogicExtClient() pb.LogicExtClient {
	conn, err := grpc.Dial("127.0.0.1:8010", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return pb.NewLogicExtClient(conn)
}

func getCtx() context.Context {
	token := "0"
	return metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
		"user_id", "3",
		"device_id", "3",
		"token", token,
		"request_id", strconv.FormatInt(time.Now().UnixNano(), 10)))
}

func TestLogicExtServer_RegisterDevice(t *testing.T) {
	resp, err := getLogicExtClient().RegisterDevice(context.TODO(),
		&pb.RegisterDeviceReq{
			Type:          1,
			Brand:         "huawei",
			Model:         "huawei P30",
			SystemVersion: "1.0.0",
			SdkVersion:    "1.0.0",
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", resp)
}

func TestLogicExtServer_SendMessageToFriend(t *testing.T) {
	resp, err := getLogicExtClient().SendMessageToFriend(getCtx(),
		&pb.SendMessageReq{
			ReceiverId: 2,
			Content:    []byte("hahaha1000"),
			SendTime:   util.UnixMilliTime(time.Now()),
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", resp)
}

func TestLogicExtServer_SendMessageToGroup(t *testing.T) {
	resp, err := getLogicExtClient().SendMessageToGroup(getCtx(),
		&pb.SendMessageReq{
			ReceiverId: 4,
			Content:    []byte("group message "),
			SendTime:   util.UnixMilliTime(time.Now()),
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", resp)
}

func TestLogicExtServer_CreateGroup(t *testing.T) {
	resp, err := getLogicExtClient().CreateGroup(getCtx(),
		&pb.CreateGroupReq{
			Name:         "10",
			Introduction: "10",
			Extra:        "10",
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", resp)
}

func TestLogicExtServer_UpdateGroup(t *testing.T) {
	resp, err := getLogicExtClient().UpdateGroup(getCtx(),
		&pb.UpdateGroupReq{
			GroupId:      2,
			Name:         "11",
			Introduction: "11",
			Extra:        "11",
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", resp)
}

func TestLogicExtServer_GetGroup(t *testing.T) {
	resp, err := getLogicExtClient().GetGroup(getCtx(),
		&pb.GetGroupReq{
			GroupId: 2,
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", resp)
}

func TestLogicExtServer_GetUserGroups(t *testing.T) {
	resp, err := getLogicExtClient().GetGroups(getCtx(), &emptypb.Empty{})
	if err != nil {
		fmt.Println(err)
		return
	}
	// todo 不能获取用户所在的超大群组
	fmt.Printf("%+v\n", resp)
}

func TestLogicExtServer_AddGroupMember(t *testing.T) {
	resp, err := getLogicExtClient().AddGroupMembers(getCtx(),
		&pb.AddGroupMembersReq{
			GroupId: 2,
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", resp)
}

func TestLogicExtServer_UpdateGroupMember(t *testing.T) {
	resp, err := getLogicExtClient().UpdateGroupMember(getCtx(),
		&pb.UpdateGroupMemberReq{
			GroupId: 2,
			UserId:  3,
			Remarks: "2",
			Extra:   "2",
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", resp)
}

func TestLogicExtServer_DeleteGroupMember(t *testing.T) {
	resp, err := getLogicExtClient().DeleteGroupMember(getCtx(),
		&pb.DeleteGroupMemberReq{
			GroupId: 10,
			UserId:  1,
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", resp)
}

func TestLogicExtServer_GetGroupMembers(t *testing.T) {
	resp, err := getLogicExtClient().GetGroupMembers(getCtx(),
		&pb.GetGroupMembersReq{
			GroupId: 2,
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", resp)
}

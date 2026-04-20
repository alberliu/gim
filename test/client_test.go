package test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gim/pkg/protocol/pb/connectpb"
	pb "gim/pkg/protocol/pb/logicpb"
)

func TestClient(t *testing.T) {
	initData()

	connect(1, 11)
	connect(1, 12)
	connect(2, 2)
	connect(3, 3)

	time.Sleep(2 * time.Second)
	fmt.Println()
	reply, err := getMessageIntClient().PushToUsers(context.TODO(), &pb.PushToUsersRequest{
		UserIds:   []uint64{1},
		Command:   connectpb.MessageCommand_MC_USER_MESSAGE,
		Content:   []byte("hello gim"),
		IsPersist: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	slog.Info("私聊发送", "MessageID", reply.MessageId)

	time.Sleep(1 * time.Second)
	fmt.Println()
	groupReply, err := getGroupIntClient().Push(context.TODO(), &pb.GroupPushRequest{
		GroupId:   1,
		Command:   connectpb.MessageCommand_MC_GROUP_MESSAGE,
		Content:   []byte("hello gim from group"),
		IsPersist: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	slog.Info("群组发送成功", "MessageID", groupReply.MessageId)

	time.Sleep(1 * time.Second)
	fmt.Println()
	_, err = getRoomIntClient().PushRoom(context.TODO(), &pb.PushRoomRequest{
		RoomId:     1,
		Command:    10000,
		Content:    []byte("hello gim from room"),
		SendTime:   time.Now().Unix(),
		IsPriority: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	slog.Info("房间发送成功")

	select {}
}

func getMessageIntClient() pb.MessageIntServiceClient {
	conn, err := grpc.NewClient(logicServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return pb.NewMessageIntServiceClient(conn)
}

func getGroupIntClient() pb.GroupIntServiceClient {
	conn, err := grpc.NewClient(logicServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return pb.NewGroupIntServiceClient(conn)
}

func getRoomIntClient() pb.RoomIntServiceClient {
	conn, err := grpc.NewClient(logicServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return pb.NewRoomIntServiceClient(conn)
}

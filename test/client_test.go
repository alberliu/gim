package test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "gim/pkg/protocol/pb/logicpb"
)

func TestClient(t *testing.T) {
	initData()

	connect(1, 11)
	connect(1, 12)
	connect(2, 2)
	connect(3, 3)

	time.Sleep(2 * time.Second)
	reply, err := getMessageIntClient().PushToUsers(context.TODO(), &pb.PushToUsersRequest{
		UserIds:   []uint64{1},
		Command:   200,
		Content:   []byte("hello gim"),
		IsPersist: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	slog.Info("私聊发送", "MessageID", reply.MessageId)

	time.Sleep(1 * time.Second)
	groupReply, err := getGroupIntClient().Push(context.TODO(), &pb.GroupPushRequest{
		GroupId:   1,
		Command:   200,
		Content:   []byte("hello gim from group"),
		IsPersist: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	slog.Info("群组发送成功", "MessageID", groupReply.MessageId)

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

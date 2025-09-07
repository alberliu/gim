package api

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "gim/pkg/protocol/pb/logicpb"
)

func getClient() pb.MessageIntServiceClient {
	conn, err := grpc.NewClient("127.0.0.1:8010", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return pb.NewMessageIntServiceClient(conn)
}

func TestMessageIntService_PushToUsers(t *testing.T) {
	reply, err := getClient().PushToUsers(context.TODO(), &pb.PushToUsersRequest{
		UserIds:   []uint64{1},
		Command:   200,
		Content:   []byte("hello gim"),
		IsPersist: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply)
}

func TestMessageIntService_PushsLocal(t *testing.T) {
	reply, err := new(MessageIntService).PushToUsers(context.TODO(), &pb.PushToUsersRequest{
		UserIds:   []uint64{1},
		Command:   100,
		Content:   []byte("hello gim3"),
		IsPersist: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply)
}

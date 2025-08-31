package room

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "gim/pkg/protocol/pb/logicpb"
)

func getIntClient() pb.RoomIntServiceClient {
	conn, err := grpc.NewClient("127.0.0.1:8010", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return pb.NewRoomIntServiceClient(conn)
}

func TestRoomIntService_PushRoom(t *testing.T) {
	reply, err := getIntClient().PushRoom(context.TODO(), &pb.PushRoomRequest{
		RoomId:     1,
		Command:    1000,
		Content:    []byte("room msg"),
		SendTime:   time.Now().Unix(),
		IsPriority: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply)
}

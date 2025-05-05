package room

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"gim/pkg/md"
	pb "gim/pkg/protocol/pb/logicpb"
)

func getExtClient() pb.RoomExtServiceClient {
	conn, err := grpc.NewClient("127.0.0.1:8010", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return pb.NewRoomExtServiceClient(conn)
}

func TestRoomExtService_PushRoom(t *testing.T) {
	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.New(map[string]string{
		md.CtxUserID:   "1",
		md.CtxDeviceID: "1",
		md.CtxToken:    "0",
	}))

	reply, err := getExtClient().PushRoom(ctx, &pb.PushRoomRequest{
		RoomId:     1,
		Code:       1000,
		Content:    []byte("room msg"),
		SendTime:   time.Now().Unix(),
		IsPersist:  false,
		IsPriority: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply)
}

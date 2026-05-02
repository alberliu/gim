package api

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "gim/pkg/gen/proto/businesspb"
	"gim/pkg/md"
)

func TestFriendExtService_Add(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		md.UserID:   "2",
		md.DeviceID: "2",
	}))

	reply, err := new(FriendExtService).Add(ctx, &pb.FriendAddRequest{
		FriendId:    1,
		Remarks:     "1号朋友",
		Description: "我是2号朋友",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply)
}

func TestFriendExtService_Agree(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		md.UserID:   "1",
		md.DeviceID: "1",
	}))

	reply, err := new(FriendExtService).Agree(ctx, &pb.FriendAgreeRequest{
		UserId:  2,
		Remarks: "2号朋友",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply)
}

func TestFriendExtService_GetFriends(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		md.UserID:   "2",
		md.DeviceID: "2",
	}))

	reply, err := new(FriendExtService).GetFriends(ctx, &emptypb.Empty{})
	if err != nil {
		t.Fatal(err)
	}
	for _, friend := range reply.Friends {
		t.Log(friend)
	}
}

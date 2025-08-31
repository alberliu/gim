package friend

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"gim/pkg/local"
	"gim/pkg/md"
	pb "gim/pkg/protocol/pb/businesspb"
)

func TestFriendExtService_Add(t *testing.T) {
	local.Init()

	ctx := metadata.NewIncomingContext(context.TODO(), metadata.New(map[string]string{
		md.CtxUserID:   "2",
		md.CtxDeviceID: "2",
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
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.New(map[string]string{
		md.CtxUserID:   "1",
		md.CtxDeviceID: "1",
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

func TestFriendExtService_SendMessage(t *testing.T) {
	local.Init()

	ctx := metadata.NewIncomingContext(context.TODO(), metadata.New(map[string]string{
		md.CtxUserID:   "2",
		md.CtxDeviceID: "2",
	}))

	reply, err := new(FriendExtService).SendMessage(ctx, &pb.SendFriendMessageRequest{
		UserId:  1,
		Content: []byte("hello im 2 2"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply)
}

func TestFriendExtService_GetFriends(t *testing.T) {
	local.Init()

	ctx := metadata.NewIncomingContext(context.TODO(), metadata.New(map[string]string{
		md.CtxUserID:   "2",
		md.CtxDeviceID: "2",
	}))

	reply, err := new(FriendExtService).GetFriends(ctx, &emptypb.Empty{})
	if err != nil {
		t.Fatal(err)
	}
	for _, friend := range reply.Friends {
		t.Log(friend)
	}
}

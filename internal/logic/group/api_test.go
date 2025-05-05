package group

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"gim/pkg/md"
	pb "gim/pkg/protocol/pb/logicpb"
)

func getExtClient() pb.GroupExtServiceClient {
	conn, err := grpc.NewClient("127.0.0.1:8010", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return pb.NewGroupExtServiceClient(conn)
}

func TestGroupExtService_Create(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.New(map[string]string{
		md.CtxUserID:   "1",
		md.CtxDeviceID: "1",
	}))

	reply, err := new(GroupExtService).Create(ctx, &pb.GroupCreateRequest{
		Name:         "群组A",
		AvatarUrl:    "",
		Introduction: "群组A的介绍",
		Extra:        "",
		MemberIds:    []uint64{2},
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(reply)
}

func TestGroupExtService_Update(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.New(map[string]string{
		md.CtxUserID:   "1",
		md.CtxDeviceID: "1",
	}))

	reply, err := new(GroupExtService).Update(ctx, &pb.GroupUpdateRequest{
		GroupId:      5,
		Name:         "群组B",
		AvatarUrl:    "",
		Introduction: "群组B的介绍",
		Extra:        "",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(reply)
}

func TestGroupExtService_Get(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.New(map[string]string{
		md.CtxUserID:   "1",
		md.CtxDeviceID: "1",
	}))

	reply, err := new(GroupExtService).Get(ctx, &pb.GroupGetRequest{GroupId: 5})
	if err != nil {
		t.Error(err)
	}
	t.Log(reply)
}

func TestGroupExtService_List(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.New(map[string]string{
		md.CtxUserID:   "1",
		md.CtxDeviceID: "1",
	}))

	reply, err := new(GroupExtService).List(ctx, &emptypb.Empty{})
	if err != nil {
		t.Error(err)
	}
	for _, group := range reply.Groups {
		t.Log(group)
	}
}

func TestGroupExtService_AddMembers(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.New(map[string]string{
		md.CtxUserID:   "1",
		md.CtxDeviceID: "1",
	}))

	_, err := new(GroupExtService).AddMembers(ctx, &pb.AddMembersRequest{
		GroupId: 5,
		UserIds: []uint64{3},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGroupExtService_UpdateMember(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.New(map[string]string{
		md.CtxUserID:   "1",
		md.CtxDeviceID: "1",
	}))

	_, err := new(GroupExtService).UpdateMember(ctx, &pb.UpdateMemberRequest{
		GroupId:    5,
		UserId:     3,
		MemberType: 2,
		Remarks:    "1",
		Extra:      "1",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGroupExtService_DeleteMember(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.New(map[string]string{
		md.CtxUserID:   "1",
		md.CtxDeviceID: "1",
	}))

	_, err := new(GroupExtService).DeleteMember(ctx, &pb.DeleteMemberRequest{
		GroupId: 5,
		UserId:  3,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGroupExtService_GetMembers(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.New(map[string]string{
		md.CtxUserID:   "1",
		md.CtxDeviceID: "1",
	}))

	reply, err := new(GroupExtService).GetMembers(ctx, &pb.GetMembersRequest{GroupId: 5})
	if err != nil {
		t.Error(err)
	}
	for _, member := range reply.Members {
		t.Log(member)
	}
}

func TestGroupExtService_SendMessage(t *testing.T) {
	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.New(map[string]string{
		md.CtxUserID:   "1",
		md.CtxDeviceID: "1",
		md.CtxToken:    "0",
	}))

	reply, err := getExtClient().SendMessage(ctx, &pb.SendGroupMessageRequest{
		GroupId: 5,
		Content: []byte("group msg hello"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply.MessageId)
}

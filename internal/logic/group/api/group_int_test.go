package api

import (
	"context"
	"testing"

	pb "gim/pkg/protocol/pb/logicpb"
)

func TestGroupExtService_Create(t *testing.T) {
	reply, err := new(GroupIntService).Create(context.TODO(), &pb.GroupCreateRequest{
		Group: &pb.Group{
			Id:           5,
			Name:         "群组B",
			AvatarUrl:    "",
			Introduction: "群组B的介绍",
			Extra:        "",
			Members:      []uint64{1, 2, 3},
		},
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(reply)
}

func TestGroupExtService_Update(t *testing.T) {
	reply, err := new(GroupIntService).Update(context.TODO(), &pb.GroupUpdateRequest{
		Group: &pb.Group{
			Id:           5,
			Name:         "群组B",
			AvatarUrl:    "",
			Introduction: "群组B的介绍",
			Extra:        "",
			Members:      []uint64{1, 2, 3},
		},
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(reply)
}

func TestGroupExtService_Get(t *testing.T) {
	reply, err := new(GroupIntService).Get(context.TODO(), &pb.GroupGetRequest{GroupId: 5})
	if err != nil {
		t.Error(err)
	}
	t.Log(reply)
}

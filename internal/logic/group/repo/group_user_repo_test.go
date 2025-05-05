package repo

import (
	"testing"

	"gim/internal/logic/group/domain"
	pb "gim/pkg/protocol/pb/logicpb"
)

func TestGroupUserDao_ListByUserId(t *testing.T) {
	groups, err := GroupUserRepo.ListByUserId(1)
	if err != nil {
		t.Fatal(err)
	}
	for _, group := range groups {
		t.Log(group)
	}
}

func TestGroupUserDao_Get(t *testing.T) {
	groupUser, err := GroupUserRepo.Get(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(groupUser)
}

func TestGroupUserDao_Delete(t *testing.T) {
	err := GroupUserRepo.Delete(1, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_groupUserRepo_Save(t *testing.T) {
	err := GroupUserRepo.Save(&domain.GroupUser{
		GroupID:    1,
		UserID:     1,
		MemberType: pb.MemberType_GMT_MEMBER,
		Remarks:    "",
		Extra:      "",
		Status:     0,
	})
	if err != nil {
		t.Fatal(err)
	}
}

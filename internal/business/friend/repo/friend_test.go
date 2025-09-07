package repo

import (
	"testing"

	"gim/internal/business/friend/domain"
)

func Test_friendRepo_Get(t *testing.T) {
	friend, err := FriendRepo.Get(1, 2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(friend)
}

func Test_friendRepo_Save(t *testing.T) {
	err := FriendRepo.Save(&domain.Friend{
		UserID:   1,
		FriendID: 2,
	})
	t.Log(err)
}

func Test_friendRepo_List(t *testing.T) {
	friends, err := FriendRepo.List(1, domain.FriendStatusAgree)
	if err != nil {
		t.Fatal(err)
	}
	for _, friend := range friends {
		t.Log(friend)
	}
}

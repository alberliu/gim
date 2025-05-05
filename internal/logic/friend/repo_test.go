package friend

import (
	"testing"
)

func Test_friendDao_Get(t *testing.T) {
	friend, err := Repo.Get(1, 2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(friend)
}

func Test_friendDao_Save(t *testing.T) {
	err := Repo.Save(&Friend{
		UserID:   1,
		FriendID: 2,
	})
	t.Log(err)
}

func Test_friendDao_List(t *testing.T) {
	friends, err := Repo.List(1, FriendStatusAgree)
	if err != nil {
		t.Fatal(err)
	}
	for _, friend := range friends {
		t.Log(friend)
	}
}

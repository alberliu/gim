package repo

import (
	"testing"

	"gim/internal/user/domain"
)

func Test_userRepo_Get(t *testing.T) {
	user, err := UserRepo.Get(10000)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}

func Test_userRepo_Save(t *testing.T) {
	err := UserRepo.Save(&domain.User{
		PhoneNumber: "1",
		Nickname:    "1",
		Sex:         1,
		AvatarUrl:   "1",
		Extra:       "1",
	})
	if err != nil {
		t.Fatal(err)
	}
}

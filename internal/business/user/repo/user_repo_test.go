package repo

import (
	"context"
	"testing"

	"gim/internal/business/user/domain"
)

func Test_userRepo_Get(t *testing.T) {
	user, err := UserRepo.Get(context.Background(), 10000)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}

func Test_userRepo_Save(t *testing.T) {
	err := UserRepo.Save(context.Background(), &domain.User{
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

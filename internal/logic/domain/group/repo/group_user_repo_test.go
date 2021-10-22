package repo

import (
	"fmt"
	"testing"
)

func TestGroupUserDao_ListByUserId(t *testing.T) {
	groups, err := GroupUserRepo.ListByUserId(1)
	fmt.Printf("%+v\n %+v\n", groups, err)
}

func TestGroupUserDao_ListGroupUser(t *testing.T) {
	users, err := GroupUserRepo.ListUser(1)
	fmt.Printf("%+v\n %+v\n", users, err)
}

func TestGroupUserDao_Get(t *testing.T) {
	fmt.Println(GroupUserRepo.Get(1, 1))
}

func TestGroupUserDao_Delete(t *testing.T) {
	fmt.Println(GroupUserRepo.Delete(1, 1))
}

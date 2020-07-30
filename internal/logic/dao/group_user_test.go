package dao

import (
	"fmt"
	"testing"
)

func TestGroupUserDao_ListByUserId(t *testing.T) {
	groups, err := GroupUserDao.ListByUserId(1)
	fmt.Printf("%+v\n %+v\n", groups, err)
}

func TestGroupUserDao_ListGroupUser(t *testing.T) {
	users, err := GroupUserDao.ListUser(1)
	fmt.Printf("%+v\n %+v\n", users, err)
}

func TestGroupUserDao_Get(t *testing.T) {
	fmt.Println(GroupUserDao.Get(1, 1))
}

func TestGroupUserDao_Add(t *testing.T) {
	fmt.Println(GroupUserDao.Add(1, 1, "1", "1"))
}

func TestGroupUserDao_Delete(t *testing.T) {
	fmt.Println(GroupUserDao.Delete(1, 1))
}

func TestGroupUserDao_Update(t *testing.T) {
	fmt.Println(GroupUserDao.Update(1, 1, "2", "2"))
}

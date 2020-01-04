package dao

import (
	"fmt"
	"testing"
)

func TestGroupUserDao_ListByUserId(t *testing.T) {
	fmt.Println(GroupUserDao.ListByUserId(1, 1))
}

func TestGroupUserDao_ListGroupUser(t *testing.T) {
	fmt.Println(GroupUserDao.ListUser(1, 1))
}

func TestGroupUserDao_Get(t *testing.T) {
	fmt.Println(GroupUserDao.Get(1, 1, 1))
}

func TestGroupUserDao_Add(t *testing.T) {
	fmt.Println(GroupUserDao.Add(1, 1, 1, "1", "1"))
}

func TestGroupUserDao_Delete(t *testing.T) {
	fmt.Println(GroupUserDao.Delete(1, 1, 1))
}

func TestGroupUserDao_Update(t *testing.T) {
	fmt.Println(GroupUserDao.Update(1, 1, 1, "2", "2"))
}

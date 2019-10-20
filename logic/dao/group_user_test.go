package dao

import (
	"fmt"
	"testing"
)

func TestGroupUserDao_ListByUserId(t *testing.T) {
	fmt.Println(GroupUserDao.ListByUserId(ctx, 1, 1))
}

func TestGroupUserDao_ListGroupUser(t *testing.T) {
	fmt.Println(GroupUserDao.ListUser(ctx, 1, 1))
}

func TestGroupUserDao_Get(t *testing.T) {
	fmt.Println(GroupUserDao.Get(ctx, 1, 1, 1))
}

func TestGroupUserDao_Add(t *testing.T) {
	fmt.Println(GroupUserDao.Add(ctx, 1, 1, 1, "1", "1"))
}

func TestGroupUserDao_Delete(t *testing.T) {
	fmt.Println(GroupUserDao.Delete(ctx, 1, 1, 1))
}

func TestGroupUserDao_Update(t *testing.T) {
	fmt.Println(GroupUserDao.Update(ctx, 1, 1, 1, "2", "2"))
}

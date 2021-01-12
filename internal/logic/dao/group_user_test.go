package dao

import (
	"fmt"
	"gim/internal/logic/model"
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
	fmt.Println(GroupUserDao.Add(model.GroupUser{
		GroupId: 1,
		UserId:  1,
		Remarks: "1",
		Extra:   "1",
		Status:  0,
	}))
}

func TestGroupUserDao_Delete(t *testing.T) {
	fmt.Println(GroupUserDao.Delete(1, 1))
}

func TestGroupUserDao_Update(t *testing.T) {
	fmt.Println(GroupUserDao.Update(model.GroupUser{
		GroupId: 1,
		UserId:  1,
		Remarks: "1",
		Extra:   "1",
		Status:  0,
	}))
}

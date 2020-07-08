package dao

import (
	"fmt"
	"gim/internal/logic/model"
	"testing"
)

func TestGroupDao_Get(t *testing.T) {
	fmt.Println(GroupDao.Get(1, 1))
}

func TestGroupDao_Add(t *testing.T) {
	group := model.Group{
		AppId:        1,
		GroupId:      1,
		Name:         "1",
		Introduction: "1",
		Type:         1,
		Extra:        "1",
	}
	fmt.Println(GroupDao.Add(group))
}

func TestGroupDao_Update(t *testing.T) {
	fmt.Println(GroupDao.Update(1, 1, "2", "2", "3"))
}

func TestGroupDao_AddUserNum(t *testing.T) {
	fmt.Println(GroupDao.AddUserNum(1, 1, -1))
}

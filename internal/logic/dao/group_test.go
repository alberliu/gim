package dao

import (
	"fmt"
	"gim/internal/logic/model"
	"testing"
)

func TestGroupDao_Get(t *testing.T) {
	group, err := GroupDao.Get(1)
	fmt.Printf("%+v\n %+v\n", group, err)
}

func TestGroupDao_Add(t *testing.T) {
	group := model.Group{
		Name:         "1",
		Introduction: "1",
		Type:         1,
		Extra:        "1",
	}
	fmt.Println(GroupDao.Add(group))
}

func TestGroupDao_Update(t *testing.T) {
	fmt.Println(GroupDao.Update(1, "2", "2", "2", ""))
}

package repo

import (
	"fmt"
	"testing"

	"gim/internal/business/domain/user/model"
)

func TestUserDao_Add(t *testing.T) {
	id, err := UserDao.Add(model.User{
		PhoneNumber: "18829291351",
		Nickname:    "Alber",
		Sex:         1,
		AvatarUrl:   "AvatarUrl",
		Extra:       "Extra",
	})
	fmt.Printf("%+v\n %+v\n ", id, err)
}

func TestUserDao_Get(t *testing.T) {
	user, err := UserDao.Get(1)
	fmt.Printf("%+v\n %+v\n ", user, err)
}

func TestUserDao_GetByIds(t *testing.T) {
	users, err := UserDao.GetByIds([]int64{1, 2, 3})
	fmt.Printf("%+v\n %+v\n ", users, err)
}

func TestUserDao_GetByPhoneNumber(t *testing.T) {
	user, err := UserDao.GetByPhoneNumber("18829291351")
	fmt.Printf("%+v\n %+v\n ", user, err)
}

func TestUserDao_Search(t *testing.T) {
	users, err := UserDao.Search("哈哈哈")
	fmt.Printf("%+v\n %+v\n ", users, err)
}

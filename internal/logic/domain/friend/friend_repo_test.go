package friend

import (
	"fmt"
	"gim/internal/logic/model"
	"testing"
)

func Test_friendDao_Get(t *testing.T) {
	friend, err := FriendDao.Get(1, 2)
	fmt.Printf("%+v \n %+v \n", friend, err)
}

func Test_friendDao_Add(t *testing.T) {
	fmt.Println(FriendDao.Add(model.Friend{
		UserId:   1,
		FriendId: 2,
	}))
}

func Test_friendDao_Update(t *testing.T) {
	fmt.Println(FriendDao.Update(model.Friend{
		UserId:   1,
		FriendId: 2,
		Remarks:  "1",
		Extra:    "1",
		Status:   1,
	}))
}

func Test_friendDao_List(t *testing.T) {
	friends, err := FriendDao.List(1, model.FriendStatusAgree)
	fmt.Printf("%+v \n %+v \n", friends, err)
}

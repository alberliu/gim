package friend

import (
	"fmt"
	"testing"
)

func Test_friendDao_Get(t *testing.T) {
	friend, err := Repo.Get(1, 2)
	fmt.Printf("%+v \n %+v \n", friend, err)
}

func Test_friendDao_Save(t *testing.T) {
	fmt.Println(Repo.Save(&Friend{
		UserId:   1,
		FriendId: 2,
	}))
}

func Test_friendDao_List(t *testing.T) {
	friends, err := Repo.List(1, FriendStatusAgree)
	fmt.Printf("%+v \n %+v \n", friends, err)
}

package service

import (
	"fmt"
	"goim/logic/db"
	"goim/logic/model"
	"goim/public/imctx"
	"testing"
)

var ctx = imctx.NewContext(db.Factoty.GetSession())

func TestFriendService_Add(t *testing.T) {
	add := model.FriendAdd{
		UserId:      1,
		UserLabel:   "alber",
		Friend:      2,
		FriendLabel: "h",
	}
	err := FriendService.Add(ctx, add)
	if err != nil {
		fmt.Println(err)
	}
}

func TestFriendService_Delete(t *testing.T) {
	err := FriendService.Delete(ctx, 1, 2)
	if err != nil {
		fmt.Println(err)
	}
}

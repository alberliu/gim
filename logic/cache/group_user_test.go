package cache

import (
	"fmt"
	"goim/logic/db"
	"goim/logic/model"
	"goim/public/imctx"
	"goim/public/util"
	"testing"
)

var ctx = imctx.NewContext(db.Factoty.GetSession())

func TestGroupUserCache_SetAll(t *testing.T) {
	var userInfos = []model.GroupUser{
		{
			AppId:   1,
			UserId:  2,
			GroupId: 1,
			Label:   "1",
			Extra:   "1",
		},
	}
	fmt.Println(GroupUserCache.SetAll(ctx, 1, 1, userInfos))
}

func TestGroupUserCache_GetAll(t *testing.T) {
	users, err := GroupUserCache.GetAll(ctx, 1, 1)
	fmt.Println(err)
	fmt.Println(util.JsonMarshal(users))
}

func TestGroupUserCache_Get(t *testing.T) {
	user, err := GroupUserCache.Get(ctx, 1, 1, 1)
	fmt.Println(err)
	fmt.Println(util.JsonMarshal(user))
}

func TestGroupUserCache_Set(t *testing.T) {
	fmt.Println(GroupUserCache.Set(1, 1, model.GroupUser{
		AppId:   1,
		UserId:  1,
		GroupId: 0,
		Label:   "2",
		Extra:   "2",
	}))
}
func TestGroupUserCache_Del(t *testing.T) {
	fmt.Println(GroupUserCache.Del(ctx, 1, 1, 1))
}

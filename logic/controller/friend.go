package controller

import (
	"goim/logic/service"
	"strconv"
)

func init() {
	g := Engine.Group("/friend")
	g.GET("/all", handler(FriendControlelr{}.All))
	g.POST("", handler(FriendControlelr{}.Add))
	g.DELETE("/:friend_id", handler(FriendControlelr{}.Delete))
}

type FriendControlelr struct{}

// Friend 好友
func (FriendControlelr) All(c *context) {
	c.response(service.FriendService.ListUserFriend(Context(), c.userId))
}

func (FriendControlelr) Add(c *context) {
	var data struct {
		FriendId int64  `json:"friend_id"`
		Label    string `json:"label"`
	}
	if c.bindJson(&data) != nil {
		return
	}
	c.response(nil, service.FriendService.Add(Context(), c.userId, data.FriendId, data.Label))
}

func (FriendControlelr) Delete(c *context) {
	friendIdStr := c.Param("friend_id")
	friendId, err := strconv.ParseInt(friendIdStr, 10, 64)
	if err != nil {
		c.badParam(err)
		return
	}
	c.response(nil, service.FriendService.Delete(Context(), c.userId, friendId))
}

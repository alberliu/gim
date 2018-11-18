package dao

import (
	"goim/logic/model"
	"goim/public/imctx"
	"goim/public/logger"
)

type friendDao struct{}

var FriendDao = new(friendDao)

// Get 获取一个朋友关系
func (*friendDao) Get(ctx *imctx.Context, userId int64, friendId int64) (*model.Friend, error) {
	var friend model.Friend
	row := ctx.Session.QueryRow(`select id,user_id,friend_id,label,create_time,update_time 
		from t_friend where user_id = ? and friend_id = ?`,
		userId, friendId)
	err := row.Scan(&friend.Id, &friend.UserId, &friend.Id, &friend.Label, &friend.CreateTime, &friend.UpdateTime)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return &friend, nil
}

// Add 插入一条朋友关系
func (*friendDao) Add(ctx *imctx.Context, userId int64, friendId int64, label string) error {
	_, err := ctx.Session.Exec("insert ignore into t_friend(user_id,friend_id,label) values(?,?,?)",
		userId, friendId, label)
	if err != nil {
		logger.Sugar.Error(err)
	}
	return err
}

// Delete 删除一条朋友关系
func (*friendDao) Delete(ctx *imctx.Context, userId, friendId int64) error {
	_, err := ctx.Session.Exec("delete from t_friend where user_id = ? and friend_id = ? ",
		userId, friendId)
	if err != nil {
		logger.Sugar.Error(err)
	}
	return err
}

// ListFriends 获取用户的朋友列表
func (*friendDao) ListUserFriend(ctx *imctx.Context, userId int64) ([]model.UserFriend, error) {
	rows, err := ctx.Session.Query(`select f.label,u.id,u.number,u.nickname,u.sex,u.avatar 
		from t_friend f left join t_user u on f.friend_id = u.id where f.user_id = ?`, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	friends := make([]model.UserFriend, 0, 5)
	for rows.Next() {
		var user model.UserFriend
		err := rows.Scan(&user.Label, &user.UserId, &user.Number, &user.Nickname, &user.Sex, &user.Avatar)
		if err != nil {
			logger.Sugar.Error(err)
			return nil, err
		}
		friends = append(friends, user)
	}
	return friends, nil
}

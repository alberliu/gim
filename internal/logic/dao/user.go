package dao

import (
	"database/sql"
	"gim/internal/logic/db"
	"gim/internal/logic/model"
	"gim/pkg/gerrors"
)

type userDao struct{}

var UserDao = new(userDao)

// Add 插入一条用户信息
func (*userDao) Add(user model.User) (int64, error) {
	result, err := db.DBCli.Exec("insert ignore into user(app_id,user_id,nickname,sex,avatar_url,extra) values(?,?,?,?,?,?)",
		user.AppId, user.UserId, user.Nickname, user.Sex, user.AvatarUrl, user.Extra)
	if err != nil {
		return 0, gerrors.WrapError(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return affected, nil
}

// Get 获取用户信息
func (*userDao) Get(appId, userId int64) (*model.User, error) {
	row := db.DBCli.QueryRow("select nickname,sex,avatar_url,extra,create_time,update_time from user where app_id = ? and user_id = ?",
		appId, userId)
	user := model.User{
		AppId:  appId,
		UserId: userId,
	}

	err := row.Scan(&user.Nickname, &user.Sex, &user.AvatarUrl, &user.Extra, &user.CreateTime, &user.UpdateTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, gerrors.WrapError(err)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

// Update 更新用户信息
func (*userDao) Update(user model.User) error {
	_, err := db.DBCli.Exec("update user set nickname = ?,sex = ?,avatar_url = ?,extra = ? where app_id = and user_id = ?",
		user.Nickname, user.Sex, user.AvatarUrl, user.Extra, user.AppId, user.UserId)
	if err != nil {
		return gerrors.WrapError(err)
	}

	return nil
}

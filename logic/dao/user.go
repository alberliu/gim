package dao

import (
	"goim/logic/model"
	"goim/public/imctx"
	"goim/public/logger"
)

type userDao struct{}

var UserDao = new(userDao)

// Add 插入一条用户信息
func (*userDao) Add(ctx *imctx.Context, user model.User) (int64, error) {
	result, err := ctx.Session.Exec("insert ignore into t_user(number,nickname,sex,avatar,password) values(?,?,?,?,?)",
		user.Number, user.Nickname, user.Sex, user.Avatar, user.Password)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}
	return id, nil
}

// Get 获取用户信息
func (*userDao) Get(ctx *imctx.Context, id int64) (*model.User, error) {
	row := ctx.Session.QueryRow("select number,nickname,password,sex,avatar from t_user where id = ?", id)
	user := new(model.User)
	err := row.Scan(&user.Number, &user.Nickname, &user.Password, &user.Sex, &user.Avatar)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return user, err
}

// GetByNumber 获取用户信息根据手机号
func (*userDao) GetByNumber(ctx *imctx.Context, number string) (*model.User, error) {
	row := ctx.Session.QueryRow("select id,number,nickname,password,sex,avatar from t_user where number = ?", number)
	user := new(model.User)
	err := row.Scan(&user.Id, &user.Number, &user.Nickname, &user.Password, &user.Sex, &user.Avatar)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return user, err
}

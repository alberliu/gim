package dao

import (
	"database/sql"
	"gim/internal/logic/db"
	"gim/internal/logic/model"
	"gim/pkg/gerrors"
)

type appDao struct{}

var AppDao = new(appDao)

// Get 获取APP信息
func (*appDao) Get(appId int64) (*model.App, error) {
	var app model.App
	err := db.DBCli.QueryRow("select id,name,private_key,create_time,update_time from app where id = ?", appId).Scan(
		&app.Id, &app.Name, &app.PrivateKey, &app.CreateTime, &app.UpdateTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, gerrors.WrapError(err)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &app, nil
}

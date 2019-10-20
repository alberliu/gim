package dao

import (
	"database/sql"
	"goim/logic/db"
	"goim/logic/model"
	"goim/public/imctx"
	"goim/public/logger"
)

type appDao struct{}

var AppDao = new(appDao)

// Get 获取APP信息
func (*appDao) Get(ctx *imctx.Context, appId int64) (*model.App, error) {
	var app model.App
	err := db.DBCli.QueryRow("select id,name,private_key,create_time,update_time from app where id = ?", appId).Scan(
		&app.Id, &app.Name, &app.PtivateKey, &app.CreateTime, &app.UpdateTime)
	if err != nil && err != sql.ErrNoRows {
		logger.Sugar.Error(err)
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &app, nil
}

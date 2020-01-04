package dao

import (
	"database/sql"
	"gim/internal/logic/db"
	"gim/internal/logic/model"
	"gim/pkg/gerrors"
	"gim/pkg/logger"
)

type groupDao struct{}

var GroupDao = new(groupDao)

// Get 获取群组信息
func (*groupDao) Get(appId, groupId int64) (*model.Group, error) {
	row := db.DBCli.QueryRow("select name,introduction,user_num,type,extra,create_time,update_time from `group` where app_id = ? and group_id = ?",
		appId, groupId)
	group := model.Group{
		AppId:   appId,
		GroupId: groupId,
	}
	err := row.Scan(&group.Name, &group.Introduction, &group.UserNum, &group.Type, &group.Extra, &group.CreateTime, &group.UpdateTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, gerrors.WrapError(err)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &group, nil
}

// Insert 插入一条群组
func (*groupDao) Add(group model.Group) (int64, error) {
	result, err := db.DBCli.Exec("insert ignore into `group`(app_id,group_id,name,introduction,type,extra) value(?,?,?,?,?,?)",
		group.AppId, group.GroupId, group.Name, group.Introduction, group.Type, group.Extra)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}
	num, err := result.RowsAffected()
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return num, nil
}

// Update 更新群组信息
func (*groupDao) Update(appId, groupId int64, name, introduction, extra string) error {
	_, err := db.DBCli.Exec("update `group` set name = ?,introduction = ?,extra = ? where app_id = ? and group_id = ?",
		name, introduction, extra, appId, groupId)
	if err != nil {
		return gerrors.WrapError(err)
	}

	return nil
}

// AddUserNum 更新群组信息
func (*groupDao) AddUserNum(appId, groupId int64, userNum int) error {
	_, err := db.DBCli.Exec("update `group` set user_num = user_num + ? where app_id = ? and group_id = ?",
		userNum, appId, groupId)
	if err != nil {
		return gerrors.WrapError(err)
	}

	return nil
}

// UpdateUserNum 更新群组群成员人数
func (*groupDao) UpdateUserNum(appId, groupId, userNum int64) error {
	_, err := db.DBCli.Exec("update `group` set user_num = user_num + ? where app_id = ? and group_id = ?",
		userNum, appId, groupId)
	if err != nil {
		return gerrors.WrapError(err)
	}

	return nil
}

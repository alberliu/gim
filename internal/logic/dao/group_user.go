package dao

import (
	"database/sql"
	"gim/internal/logic/db"
	"gim/internal/logic/model"
	"gim/pkg/gerrors"
)

type groupUserDao struct{}

var GroupUserDao = new(groupUserDao)

// ListByUser 获取用户加入的群组信息
func (*groupUserDao) ListByUserId(appId, userId int64) ([]model.Group, error) {
	rows, err := db.DBCli.Query(
		"select g.group_id,g.name,g.introduction,g.user_num,g.type,g.extra,g.create_time,g.update_time "+
			"from group_user u "+
			"left join `group` g on u.app_id = g.app_id and u.group_id = g.group_id "+
			"where u.app_id = ? and u.user_id = ?",
		appId, userId)
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	var groups []model.Group
	var group model.Group
	for rows.Next() {
		err := rows.Scan(&group.GroupId, &group.Name, &group.Introduction, &group.UserNum, &group.Type, &group.Extra, &group.CreateTime, &group.UpdateTime)
		if err != nil {
			return nil, gerrors.WrapError(err)
		}
		groups = append(groups, group)
	}
	return groups, nil
}

// ListGroupUser 获取群组用户信息
func (*groupUserDao) ListUser(appId, groupId int64) ([]model.GroupUser, error) {
	rows, err := db.DBCli.Query(`
		select user_id,label,extra,create_time,update_time 
		from group_user
		where app_id = ? and group_id = ?`, appId, groupId)
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	groupUsers := make([]model.GroupUser, 0, 5)
	for rows.Next() {
		var groupUser = model.GroupUser{
			AppId:   appId,
			GroupId: groupId,
		}
		err := rows.Scan(&groupUser.UserId, &groupUser.Label, &groupUser.Extra, &groupUser.CreateTime, &groupUser.UpdateTime)
		if err != nil {
			return nil, gerrors.WrapError(err)
		}
		groupUsers = append(groupUsers, groupUser)
	}
	return groupUsers, nil
}

// GetGroupUser 获取群组用户信息,用户不存在返回nil
func (*groupUserDao) Get(appId, groupId, userId int64) (*model.GroupUser, error) {
	var groupUser = model.GroupUser{
		AppId:   appId,
		GroupId: groupId,
		UserId:  userId,
	}
	err := db.DBCli.QueryRow("select label,extra from group_user where app_id = ? and group_id = ? and user_id = ?",
		appId, groupId, userId).
		Scan(&groupUser.Label, &groupUser.Extra)
	if err != nil && err != sql.ErrNoRows {
		return nil, gerrors.WrapError(err)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &groupUser, nil
}

// Add 将用户添加到群组
func (*groupUserDao) Add(appId, groupId, userId int64, label, extra string) error {
	_, err := db.DBCli.Exec("insert ignore into group_user(app_id,group_id,user_id,label,extra) values(?,?,?,?,?)",
		appId, groupId, userId, label, extra)
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Delete 将用户从群组删除
func (d *groupUserDao) Delete(appId int64, groupId int64, userId int64) error {
	_, err := db.DBCli.Exec("delete from group_user where app_id = ? and group_id = ? and user_id = ?",
		appId, groupId, userId)
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Update 更新用户群组信息
func (*groupUserDao) Update(appId, groupId, userId int64, label string, extra string) error {
	_, err := db.DBCli.Exec("update group_user set label = ?,extra = ? where app_id = ? and group_id = ? and user_id = ?",
		label, extra, appId, groupId, userId)
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

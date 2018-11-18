package dao

import (
	"goim/logic/model"
	"goim/public/imctx"
	"goim/public/logger"
)

type groupUserDao struct{}

var GroupUserDao = new(groupUserDao)

func (*groupUserDao) Get(ctx *imctx.Context, id int64) (*model.Group, error) {
	row := ctx.Session.QueryRow("select id,name from t_group where id = ?", id)
	var group model.Group
	err := row.Scan(&group.Id, &group.Name)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return &group, nil
}

// ListGroupUser 获取群组用户信息
func (*groupUserDao) ListGroupUser(ctx *imctx.Context, id int64) ([]model.GroupUser, error) {
	sql := `select g.label,u.id,u.number,u.nickname,u.sex,u.avatar from t_group_user g left join t_user u on g.user_id = u.id where group_id = ?`
	rows, err := ctx.Session.Query(sql, id)
	if err != nil {
		return nil, err
	}
	groupUsers := make([]model.GroupUser, 0, 5)
	for rows.Next() {
		var groupUser model.GroupUser
		err := rows.Scan(&groupUser.Label, &groupUser.UserId, &groupUser.Number, &groupUser.Name,
			&groupUser.Sex, &groupUser.Img)
		if err != nil {
			logger.Sugar.Error(err)
			return nil, err
		}
		groupUsers = append(groupUsers, groupUser)
	}
	return groupUsers, nil
}

// ListGroupUserId 获取群组用户id列表
func (*groupUserDao) ListGroupUserId(ctx *imctx.Context, id int) ([]int, error) {
	rows, err := ctx.Session.Query("select user_id t_group_user where group_id = ?", id)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	userIds := make([]int, 0, 5)
	for rows.Next() {
		var userId int
		err := rows.Scan(&userId)
		if err != nil {
			logger.Sugar.Error(err)
			return nil, err
		}
		userIds = append(userIds, userId)
	}
	return userIds, nil
}

// ListByUser 获取用户群组id列表
func (*groupUserDao) ListbyUserId(ctx *imctx.Context, userId int64) ([]int64, error) {
	rows, err := ctx.Session.Query("select group_id from t_group_user where user_id = ?", userId)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	var ids []int64
	var id int64
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			logger.Sugar.Error(err)
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// Add 将用户添加到群组
func (*groupUserDao) Add(ctx *imctx.Context, groupId int64, userId int64) error {
	_, err := ctx.Session.Exec("insert ignore into t_group_user(group_id,user_id) values(?,?)",
		groupId, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// Delete 将用户从群组删除
func (d *groupUserDao) Delete(ctx *imctx.Context, groupId int64, userId int64) error {
	_, err := ctx.Session.Exec("delete from t_group_user where group_id = ? and user_id = ?",
		groupId, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// UpdateLabel 更新用户群组备注
func (*groupUserDao) UpdateLabel(ctx *imctx.Context, groupId int64, userId int64, label string) error {
	_, err := ctx.Session.Exec("update t_group_user set label = ? where group_id = ? and user_id = ?",
		label, groupId, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// UserInGroup 用户是否在群组中
func (*groupUserDao) UserInGroup(ctx *imctx.Context, groupId int64, userId int64) (bool, error) {
	var count int
	err := ctx.Session.QueryRow("select count(*) from t_group_user where group_id = ? and user_id = ?",
		groupId, userId).
		Scan(&count)
	if err != nil {
		logger.Sugar.Error(err)
		return false, err
	}
	if count == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

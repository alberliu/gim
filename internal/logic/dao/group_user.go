package dao

import (
	"gim/internal/logic/model"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"time"

	"github.com/jinzhu/gorm"
)

type groupUserDao struct{}

var GroupUserDao = new(groupUserDao)

// ListByUser 获取用户加入的群组信息
func (*groupUserDao) ListByUserId(userId int64) ([]model.Group, error) {
	var groups []model.Group
	err := db.DB.Select("g.id,g.name,g.introduction,g.user_num,g.type,g.extra,g.create_time,g.update_time").
		Table("group_user u").
		Joins("join `group` g on u.group_id = g.id").
		Where("u.user_id = ?", userId).
		Find(&groups).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return groups, nil
}

// ListGroupUser 获取群组用户信息
func (*groupUserDao) ListUser(groupId int64) ([]model.GroupUser, error) {
	var groupUsers []model.GroupUser
	err := db.DB.Find(&groupUsers, " group_id = ?", groupId).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return groupUsers, nil
}

// GetGroupUser 获取群组用户信息,用户不存在返回nil
func (*groupUserDao) Get(groupId, userId int64) (*model.GroupUser, error) {
	var groupUser model.GroupUser
	err := db.DB.First(&groupUser, "group_id = ? and user_id = ?", groupId, userId).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, gerrors.WrapError(err)
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &groupUser, nil
}

// Add 将用户添加到群组
func (*groupUserDao) Add(groupId, userId int64, remarks, extra string) error {
	err := db.DB.Create(&model.GroupUser{
		GroupId:    groupId,
		UserId:     userId,
		Remarks:    remarks,
		Extra:      extra,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Delete 将用户从群组删除
func (d *groupUserDao) Delete(groupId int64, userId int64) error {
	err := db.DB.Exec("delete from group_user where group_id = ? and user_id = ?",
		groupId, userId).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Update 更新用户群组信息
func (*groupUserDao) Update(groupId, userId int64, remarks, extra string) error {
	err := db.DB.Exec("update group_user set remarks = ?,extra = ? where group_id = ? and user_id = ?",
		remarks, extra, groupId, userId).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

package repo

import (
	"errors"
	"gim/internal/logic/domain/group/entity"
	"gim/pkg/db"
	"gim/pkg/gerrors"

	"github.com/jinzhu/gorm"
)

type groupUserRepo struct{}

var GroupUserRepo = new(groupUserRepo)

// ListByUserId 获取用户加入的群组信息
func (*groupUserRepo) ListByUserId(userId int64) ([]entity.Group, error) {
	var groups []entity.Group
	err := db.DB.Select("g.id,g.name,g.avatar_url,g.introduction,g.user_num,g.extra,g.create_time,g.update_time").
		Table("group_user u").
		Joins("join `group` g on u.group_id = g.id").
		Where("u.user_id = ?", userId).
		Find(&groups).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return groups, nil
}

// ListUser 获取群组用户信息
func (*groupUserRepo) ListUser(groupId int64) ([]entity.GroupUser, error) {
	var groupUsers []entity.GroupUser
	err := db.DB.Find(&groupUsers, "group_id = ?", groupId).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return groupUsers, nil
}

// Get 获取群组用户信息,用户不存在返回nil
func (*groupUserRepo) Get(groupId, userId int64) (*entity.GroupUser, error) {
	var groupUser entity.GroupUser
	err := db.DB.First(&groupUser, "group_id = ? and user_id = ?", groupId, userId).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.WrapError(err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &groupUser, nil
}

// Save 将用户添加到群组
func (*groupUserRepo) Save(groupUser *entity.GroupUser) error {
	err := db.DB.Save(&groupUser).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Delete 将用户从群组删除
func (d *groupUserRepo) Delete(groupId int64, userId int64) error {
	err := db.DB.Exec("delete from group_user where group_id = ? and user_id = ?",
		groupId, userId).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

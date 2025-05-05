package repo

import (
	"errors"

	"gorm.io/gorm"

	"gim/internal/logic/group/domain"
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

type groupUserRepo struct{}

var GroupUserRepo = new(groupUserRepo)

// ListByUserId 获取用户加入的群组信息
func (*groupUserRepo) ListByUserId(userID uint64) ([]domain.Group, error) {
	var groupUsers []domain.GroupUser
	err := db.DB.Find(&groupUsers, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	groupIds := make([]uint64, 0, len(groupUsers))
	for _, groupUser := range groupUsers {
		groupIds = append(groupIds, groupUser.GroupID)
	}
	var groups []domain.Group
	err = db.DB.Find(&groups, "id in (?)", groupIds).Error
	return groups, err
}

// Get 获取群组成员信息
func (*groupUserRepo) Get(groupId, userId uint64) (*domain.GroupUser, error) {
	var groupUser domain.GroupUser
	err := db.DB.First(&groupUser, "group_id = ? and user_id = ?", groupId, userId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.ErrGroupUserNotFound
	}
	return &groupUser, err
}

func (*groupUserRepo) BatchCreate(groupUsers []domain.GroupUser) error {
	if len(groupUsers) == 0 {
		return nil
	}
	return db.DB.Create(&groupUsers).Error
}

// Save 将用户添加到群组
func (*groupUserRepo) Save(groupUser *domain.GroupUser) error {
	return db.DB.Where("group_id = ? and user_id = ?", groupUser.GroupID, groupUser.UserID).Save(&groupUser).Error
}

// Delete 将用户从群组删除
func (d *groupUserRepo) Delete(groupId, userId uint64) error {
	return db.DB.Delete(&domain.GroupUser{}, "group_id = ? and user_id = ?", groupId, userId).Error
}

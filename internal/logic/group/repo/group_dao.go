package repo

import (
	"errors"

	"gorm.io/gorm"

	"gim/internal/logic/group/domain"
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

type groupDao struct{}

var GroupDao = new(groupDao)

// Get 获取群组信息
func (*groupDao) Get(groupID uint64) (*domain.Group, error) {
	var group domain.Group
	err := db.DB.Preload("Members").First(&group, "id = ?", groupID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.ErrGroupNotFound
	}
	return &group, err
}

// Save 插入一条群组
func (*groupDao) Save(group *domain.Group) error {
	return db.DB.Save(&group).Error
}

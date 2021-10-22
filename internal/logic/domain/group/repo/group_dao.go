package repo

import (
	"gim/internal/logic/domain/group/model"
	"gim/pkg/db"
	"gim/pkg/gerrors"

	"github.com/jinzhu/gorm"
)

type groupDao struct{}

var GroupDao = new(groupDao)

// Get 获取群组信息
func (*groupDao) Get(groupId int64) (*model.Group, error) {
	var group = model.Group{Id: groupId}
	err := db.DB.First(&group).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, gerrors.WrapError(err)
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &group, nil
}

// Save 插入一条群组
func (*groupDao) Save(group *model.Group) error {
	err := db.DB.Save(&group).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

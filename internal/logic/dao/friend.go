package dao

import (
	"gim/internal/logic/model"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"time"

	"github.com/jinzhu/gorm"
)

type friendDao struct{}

var FriendDao = new(friendDao)

// Get 获取好友
func (*friendDao) Get(userId, friendId int64) (*model.Friend, error) {
	friend := model.Friend{}
	err := db.DB.First(&friend, "user_id = ? and friend_id = ?", userId, friendId).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &friend, nil
}

// Add 添加好友
func (*friendDao) Add(friend model.Friend) error {
	friend.CreateTime = time.Now()
	friend.UpdateTime = time.Now()
	return gerrors.WrapError(db.DB.Create(&friend).Error)
}

// Update 更新好友
func (*friendDao) Update(friend model.Friend) error {
	err := db.DB.Model(&friend).Where("user_id = ? and friend_id = ?", friend.UserId, friend.FriendId).
		Updates(
			map[string]interface{}{
				"remarks": friend.Remarks,
				"extra":   friend.Extra,
			},
		).Error
	return gerrors.WrapError(err)
}

// UpdateStatus 更新好友状态
func (*friendDao) UpdateStatus(userId, friendId int64, status int) error {
	err := db.DB.Model(&model.Friend{}).Where("user_id = ? and friend_id = ?", userId, friendId).
		Updates(map[string]interface{}{
			"status": status,
		}).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// List 获取好友列表
func (*friendDao) List(userId int64, status int) ([]model.Friend, error) {
	var friends []model.Friend
	err := db.DB.Where("user_id = ? and status = ?", userId, status).Find(&friends).Error
	return friends, gerrors.WrapError(err)
}

package friend

import (
	"errors"

	"gorm.io/gorm"

	"gim/pkg/db"
	"gim/pkg/gerrors"
)

type repo struct{}

var Repo = new(repo)

// Get 获取好友
func (*repo) Get(userId, friendId uint64) (*Friend, error) {
	friend := Friend{}
	err := db.DB.First(&friend, "user_id = ? and friend_id = ?", userId, friendId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.ErrFriendNotFound
	}
	return &friend, err
}

// Create 添加好友
func (*repo) Create(friend *Friend) error {
	return db.DB.Create(friend).Error
}

// Save 添加好友
func (*repo) Save(friend *Friend) error {
	return db.DB.Where("user_id = ? and friend_id = ?", friend.UserID, friend.FriendID).Save(friend).Error
}

// List 获取好友列表
func (*repo) List(userId uint64, status int) ([]Friend, error) {
	var friends []Friend
	err := db.DB.Where("user_id = ? and status = ?", userId, status).Find(&friends).Error
	return friends, err
}

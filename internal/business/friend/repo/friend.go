package repo

import (
	"errors"

	"gorm.io/gorm"

	"gim/internal/business/friend/domain"
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

var FriendRepo = new(friendRepo)

type friendRepo struct{}

// Get 获取好友
func (*friendRepo) Get(userId, friendId uint64) (*domain.Friend, error) {
	friend := domain.Friend{}
	err := db.DB.First(&friend, "user_id = ? and friend_id = ?", userId, friendId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.ErrFriendNotFound
	}
	return &friend, err
}

// Create 添加好友
func (*friendRepo) Create(friend *domain.Friend) error {
	return db.DB.Create(friend).Error
}

// Save 添加好友
func (*friendRepo) Save(friend *domain.Friend) error {
	return db.DB.Where("user_id = ? and friend_id = ?", friend.UserID, friend.FriendID).Save(friend).Error
}

// List 获取好友列表
func (*friendRepo) List(userId uint64, status int) ([]domain.Friend, error) {
	var friends []domain.Friend
	err := db.DB.Where("user_id = ? and status = ?", userId, status).Find(&friends).Error
	return friends, err
}

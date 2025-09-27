package repo

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gim/internal/business/friend/domain"
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

var FriendRepo = new(friendRepo)

type friendRepo struct{}

// Get 获取好友
func (*friendRepo) Get(ctx context.Context, userId, friendId uint64) (*domain.Friend, error) {
	friend := domain.Friend{}
	err := db.DB.WithContext(ctx).First(&friend, "user_id = ? and friend_id = ?", userId, friendId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.ErrFriendNotFound
	}
	return &friend, err
}

// Create 添加好友
func (*friendRepo) Create(ctx context.Context, friend *domain.Friend) error {
	return db.DB.WithContext(ctx).Create(friend).Error
}

// Save 添加好友
func (*friendRepo) Save(ctx context.Context, friend *domain.Friend) error {
	return db.DB.WithContext(ctx).Where("user_id = ? and friend_id = ?", friend.UserID, friend.FriendID).Save(friend).Error
}

// List 获取好友列表
func (*friendRepo) List(ctx context.Context, userId uint64, status int) ([]domain.Friend, error) {
	var friends []domain.Friend
	err := db.DB.WithContext(ctx).Where("user_id = ? and status = ?", userId, status).Find(&friends).Error
	return friends, err
}

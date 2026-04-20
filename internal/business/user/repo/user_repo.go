package repo

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gim/internal/business/user/domain"
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

type userRepo struct{}

var UserRepo = new(userRepo)

// Get 获取单个用户
func (*userRepo) Get(ctx context.Context, userID uint64) (*domain.User, error) {
	var user domain.User
	err := db.DB.WithContext(ctx).First(&user, "id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.ErrUserNotFound
	}
	return &user, err
}

func (*userRepo) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	var user domain.User
	err := db.DB.WithContext(ctx).First(&user, "phone_number = ?", phoneNumber).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.ErrUserNotFound
	}
	return &user, err
}

// GetByIDs 获取多个用户
func (*userRepo) GetByIDs(ctx context.Context, userIDs []uint64) ([]domain.User, error) {
	var users []domain.User
	err := db.DB.WithContext(ctx).Find(&users, "id in (?)", userIDs).Error
	return users, err
}

// Search 搜索用户
func (*userRepo) Search(ctx context.Context, key string) ([]domain.User, error) {
	var users []domain.User
	key = "%" + key + "%"
	err := db.DB.WithContext(ctx).Where("phone_number like ? or nickname like ?", key, key).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Save 保存用户
func (*userRepo) Save(ctx context.Context, user *domain.User) error {
	return db.DB.WithContext(ctx).Save(user).Error
}

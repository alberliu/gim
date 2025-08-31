package repo

import (
	"errors"

	"gorm.io/gorm"

	"gim/internal/business/user/domain"
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

type userRepo struct{}

var UserRepo = new(userRepo)

// Get 获取单个用户
func (*userRepo) Get(userID uint64) (*domain.User, error) {
	var user domain.User
	err := db.DB.First(&user, "id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.ErrUserNotFound
	}
	return &user, err
}

func (*userRepo) GetByPhoneNumber(phoneNumber string) (*domain.User, error) {
	var user domain.User
	err := db.DB.First(&user, "phone_number = ?", phoneNumber).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.ErrUserNotFound
	}
	return &user, err
}

// GetByIDs 获取多个用户
func (*userRepo) GetByIDs(userIDs []uint64) ([]domain.User, error) {
	var users []domain.User
	err := db.DB.Find(&users, "id in (?)", userIDs).Error
	return users, err
}

// Search 搜索用户
func (*userRepo) Search(key string) ([]domain.User, error) {
	var users []domain.User
	key = "%" + key + "%"
	err := db.DB.Where("phone_number like ? or nickname like ?", key, key).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Save 保存用户
func (*userRepo) Save(user *domain.User) error {
	return db.DB.Save(user).Error
}

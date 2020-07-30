package dao

import (
	"gim/internal/user/model"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"time"

	"github.com/jinzhu/gorm"
)

type userDao struct{}

var UserDao = new(userDao)

// Add 插入一条用户信息
func (*userDao) Add(user model.User) (int64, error) {
	user.CreateTime = time.Now()
	user.UpdateTime = time.Now()
	err := db.DB.Create(&user).Error
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return user.Id, nil
}

// Get 获取用户信息
func (*userDao) Get(userId int64) (*model.User, error) {
	var user = model.User{Id: userId}
	err := db.DB.First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, gerrors.WrapError(err)
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// GetByPhoneNumber 根据手机号获取用户信息
func (*userDao) GetByPhoneNumber(phoneNumber string) (*model.User, error) {
	var user model.User
	err := db.DB.First(&user, "phone_number = ?", phoneNumber).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, gerrors.WrapError(err)
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// Get 获取用户信息
func (*userDao) GetByIds(userIds []int64) ([]model.User, error) {
	var users []model.User
	err := db.DB.Find(&users, "id in (?)", userIds).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return users, err
}

// Update 更新用户信息
func (*userDao) Update(user model.User) error {
	err := db.DB.Model(&user).Updates(map[string]interface{}{
		"nickname":   user.Nickname,
		"sex":        user.Sex,
		"avatar_url": user.AvatarUrl,
		"extra":      user.Extra,
	}).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

package repo

import (
	"gim/internal/business/domain/user/model"
)

type userRepo struct{}

var UserRepo = new(userRepo)

// Get 获取单个用户
func (*userRepo) Get(userId int64) (*model.User, error) {
	user, err := UserCache.Get(userId)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	user, err = UserDao.Get(userId)
	if err != nil {
		return nil, err
	}

	if user != nil {
		err = UserCache.Set(*user)
		if err != nil {
			return nil, err
		}
	}
	return user, err
}

func (*userRepo) GetByPhoneNumber(phoneNumber string) (*model.User, error) {
	return UserDao.GetByPhoneNumber(phoneNumber)
}

// GetByIds 获取多个用户
func (*userRepo) GetByIds(userIds []int64) ([]model.User, error) {
	return UserDao.GetByIds(userIds)
}

// Search 搜索用户
func (*userRepo) Search(key string) ([]model.User, error) {
	return UserDao.Search(key)
}

// Save 保存用户
func (*userRepo) Save(user *model.User) error {
	userId := user.Id
	err := UserDao.Save(user)
	if err != nil {
		return err
	}

	if userId != 0 {
		err = UserCache.Del(user.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

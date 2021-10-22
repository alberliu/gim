package repo

import (
	"gim/internal/business/domain/user/model"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	UserKey    = "user:"
	UserExpire = 2 * time.Hour
)

type userCache struct{}

var UserCache = new(userCache)

// Get 获取用户缓存
func (c *userCache) Get(userId int64) (*model.User, error) {
	var user model.User
	err := db.RedisUtil.Get(UserKey+strconv.FormatInt(userId, 10), &user)
	if err != nil && err != redis.Nil {
		return nil, gerrors.WrapError(err)
	}
	if err == redis.Nil {
		return nil, nil
	}
	return &user, nil
}

// Set 设置用户缓存
func (c *userCache) Set(user model.User) error {
	err := db.RedisUtil.Set(UserKey+strconv.FormatInt(user.Id, 10), user, UserExpire)
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Del 删除用户缓存
func (c *userCache) Del(userId int64) error {
	_, err := db.RedisCli.Del(UserKey + strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

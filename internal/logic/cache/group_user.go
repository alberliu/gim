package cache

import (
	"gim/internal/logic/model"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	GroupUserKey = "group_user:"
	GroupUserExp = 2 * time.Hour
)

type groupUserCache struct{}

var GroupUserCache = new(groupUserCache)

// Set 保存群组所有用户的信息
func (c *groupUserCache) Set(groupId int64, userInfos []model.GroupUser) error {
	err := RedisUtil.Set(GroupUserKey+strconv.FormatInt(groupId, 10), userInfos, GroupUserExp)
	return gerrors.WrapError(err)
}

// GetAll 获取群组的所有用户，如果缓存里面没有，返回nil
func (c *groupUserCache) Get(groupId int64) ([]model.GroupUser, error) {
	var users []model.GroupUser
	err := RedisUtil.Get(GroupUserKey+strconv.FormatInt(groupId, 10), &users)
	if err != nil && err != redis.Nil {
		return nil, gerrors.WrapError(err)
	}
	if err == redis.Nil {
		return nil, nil
	}
	return users, nil
}

// Del 删除缓存
func (c *groupUserCache) Del(groupId int64) error {
	_, err := db.RedisCli.Del(GroupUserKey + strconv.FormatInt(groupId, 10)).Result()
	return gerrors.WrapError(err)
}

package repo

import (
	"gim/internal/logic/domain/group/entity"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const GroupKey = "group:"

type groupCache struct{}

var GroupCache = new(groupCache)

// Get 获取群组缓存
func (c *groupCache) Get(groupId int64) (*entity.Group, error) {
	var user entity.Group
	err := db.RedisUtil.Get(GroupKey+strconv.FormatInt(groupId, 10), &user)
	if err != nil && err != redis.Nil {
		return nil, gerrors.WrapError(err)
	}
	if err == redis.Nil {
		return nil, nil
	}
	return &user, nil
}

// Set 设置群组缓存
func (c *groupCache) Set(group *entity.Group) error {
	err := db.RedisUtil.Set(GroupKey+strconv.FormatInt(group.Id, 10), group, 24*time.Hour)
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Del 删除群组缓存
func (c *groupCache) Del(groupId int64) error {
	_, err := db.RedisCli.Del(GroupKey + strconv.FormatInt(groupId, 10)).Result()
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

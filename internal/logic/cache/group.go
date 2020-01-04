package cache

import (
	"gim/internal/logic/db"
	"gim/internal/logic/model"
	"gim/pkg/gerrors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	GroupKey    = "group:"
	GroupExpire = 2 * time.Hour
)

type groupCache struct{}

var GroupCache = new(groupCache)

func (*groupCache) Key(appId, groupId int64) string {
	return GroupKey + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(groupId, 10)
}

// Get 获取群组缓存
func (c *groupCache) Get(appId, groupId int64) (*model.Group, error) {
	var user model.Group
	err := get(c.Key(appId, groupId), &user)
	if err != nil && err != redis.Nil {
		return nil, gerrors.WrapError(err)
	}
	if err == redis.Nil {
		return nil, nil
	}
	return &user, nil
}

// Set 设置群组缓存
func (c *groupCache) Set(group *model.Group) error {
	err := set(c.Key(group.AppId, group.GroupId), group, GroupExpire)
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Del 删除群组缓存
func (c *groupCache) Del(appId, groupId int64) error {
	_, err := db.RedisCli.Del(c.Key(appId, groupId)).Result()
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

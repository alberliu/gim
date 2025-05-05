package repo

import (
	"strconv"
	"time"

	"gim/internal/logic/group/domain"
	"gim/pkg/db"
)

const GroupKey = "group:"

type groupCache struct{}

var GroupCache = new(groupCache)

// Get 获取群组缓存
func (c *groupCache) Get(groupId uint64) (*domain.Group, error) {
	var user domain.Group
	err := db.RedisUtil.Get(GroupKey+strconv.FormatUint(groupId, 10), &user)
	return &user, err
}

// Set 设置群组缓存
func (c *groupCache) Set(group *domain.Group) error {
	return db.RedisUtil.Set(GroupKey+strconv.FormatUint(group.ID, 10), group, 24*time.Hour)
}

// Delete 删除群组缓存
func (c *groupCache) Delete(groupId uint64) error {
	_, err := db.RedisCli.Del(GroupKey + strconv.FormatUint(groupId, 10)).Result()
	return err
}

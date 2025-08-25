package repo

import (
	"fmt"
	"time"

	"gim/internal/logic/group/domain"
	"gim/pkg/db"
)

const GroupKey = "group:%d"

type groupCache struct{}

var GroupCache = new(groupCache)

// Get 获取群组缓存
func (c *groupCache) Get(groupId uint64) (*domain.Group, error) {
	var user domain.Group
	err := db.RedisCli.GetObject(fmt.Sprintf(GroupKey, groupId), &user)
	return &user, err
}

// Set 设置群组缓存
func (c *groupCache) Set(group *domain.Group) error {
	return db.RedisCli.SetObject(fmt.Sprintf(GroupKey, group.ID), group, 24*time.Hour)
}

// Delete 删除群组缓存
func (c *groupCache) Delete(groupId uint64) error {
	_, err := db.RedisCli.Del(fmt.Sprintf(GroupKey, groupId)).Result()
	return err
}

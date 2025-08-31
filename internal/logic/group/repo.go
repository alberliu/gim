package group

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"

	"gim/pkg/db"
	"gim/pkg/gerrors"
)

const GroupKey = "group:%d"

type repo struct{}

var Repo = new(repo)

// Get 获取群组信息
func (*repo) Get(groupId uint64) (*Group, error) {
	key := fmt.Sprintf(GroupKey, groupId)
	var group Group
	err := db.RedisCli.GetAny(key, &group)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if err == nil {
		return &group, nil
	}

	err = db.DB.First(&group, "id = ?", groupId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.ErrGroupNotFound
	}
	if err != nil {
		return nil, err
	}

	err = db.RedisCli.SetAny(key, &group, 24*time.Hour)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (*repo) Create(group *Group) error {
	return db.DB.Create(group).Error
}

// Save 修改群组信息
func (*repo) Save(group *Group) error {
	err := db.DB.Save(group).Error
	if err != nil {
		return err
	}

	key := fmt.Sprintf(GroupKey, group.ID)
	return db.RedisCli.Del(key).Err()
}

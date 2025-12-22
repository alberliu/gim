package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gim/internal/logic/group/domain"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"gim/pkg/uredis"
)

const GroupKey = "group:%d"

var GroupRepo = new(groupRepo)

type groupRepo struct{}

// Get 获取群组信息
func (*groupRepo) Get(ctx context.Context, groupID uint64) (*domain.Group, error) {
	key := fmt.Sprintf(GroupKey, groupID)
	return uredis.Get(db.RedisCli, ctx, key, 24*time.Hour, func() (*domain.Group, error) {
		group, err := gorm.G[domain.Group](db.DB).Where("id = ?", groupID).First(ctx)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gerrors.ErrGroupNotFound
		}
		if err != nil {
			return nil, err
		}
		return &group, nil
	})
}

func (*groupRepo) Create(ctx context.Context, group *domain.Group) error {
	return gorm.G[domain.Group](db.DB).Create(ctx, group)
}

// Save 修改群组信息
func (*groupRepo) Save(ctx context.Context, group *domain.Group) error {
	err := db.DB.WithContext(ctx).Save(group).Error
	if err != nil {
		return err
	}

	key := fmt.Sprintf(GroupKey, group.ID)
	return db.RedisCli.Del(ctx, key).Err()
}

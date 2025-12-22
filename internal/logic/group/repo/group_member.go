package repo

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gim/internal/logic/group/domain"
	"gim/pkg/db"
	"gim/pkg/uredis"
)

const GroupMemberKey = "groupMember:%d"

type groupMemberRepo struct{}

var GroupMemberRepo groupMemberRepo

func (r *groupMemberRepo) ListByGroupID(ctx context.Context, groupID uint64) ([]domain.GroupMember, error) {
	key := fmt.Sprintf(GroupMemberKey, groupID)
	result, err := uredis.Get(db.RedisCli, ctx, key, 24*time.Hour, func() (*[]domain.GroupMember, error) {
		members, err := gorm.G[domain.GroupMember](db.DB).Where("group_id = ?", groupID).Find(ctx)
		if err != nil {
			return nil, err
		}
		return &members, nil
	})
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (r *groupMemberRepo) BatchCreate(ctx context.Context, groupID uint64, members []domain.GroupMember) error {
	err := gorm.G[domain.GroupMember](db.DB).CreateInBatches(ctx, &members, 10)
	if err != nil {
		return err
	}
	return r.deleteCache(ctx, groupID)
}

func (r *groupMemberRepo) Update(ctx context.Context, member *domain.GroupMember) error {
	err := db.DB.WithContext(ctx).Save(member).Error
	if err != nil {
		return err
	}
	return r.deleteCache(ctx, member.GroupID)
}

func (r *groupMemberRepo) BatchDelete(ctx context.Context, groupID uint64, userIDs []uint64) error {
	err := db.DB.WithContext(ctx).Where("group_id = ? and user_id in ?", groupID, userIDs).Delete(&domain.GroupMember{}).Error
	if err != nil {
		return err
	}
	return r.deleteCache(ctx, groupID)
}

func (r *groupMemberRepo) deleteCache(ctx context.Context, groupID uint64) error {
	key := fmt.Sprintf(GroupMemberKey, groupID)
	return db.RedisCli.Del(ctx, key).Err()
}

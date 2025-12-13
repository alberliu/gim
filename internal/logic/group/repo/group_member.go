package repo

import (
	"context"

	"gorm.io/gorm"

	"gim/internal/logic/group/domain"
	"gim/pkg/db"
)

type groupMemberRepo struct{}

var GroupMemberRepo groupMemberRepo

func (r *groupMemberRepo) ListByGroupID(ctx context.Context, groupID uint64) ([]domain.GroupMember, error) {
	var members []domain.GroupMember
	err := db.DB.WithContext(ctx).Where("group_id = ?", groupID).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *groupMemberRepo) Get(ctx context.Context, groupID, userID uint64) (*domain.GroupMember, error) {
	return gorm.G[*domain.GroupMember](db.DB).Where("group_id = ? and user_id = ?", groupID, userID).First(ctx)
}

func (r *groupMemberRepo) BatchCreate(ctx context.Context, members []domain.GroupMember) error {
	return gorm.G[domain.GroupMember](db.DB).CreateInBatches(ctx, &members, 10)
}

func (r *groupMemberRepo) Update(ctx context.Context, member *domain.GroupMember) error {
	return db.DB.WithContext(ctx).Save(member).Error
}

func (r *groupMemberRepo) BatchDelete(ctx context.Context, groupID uint64, userIDs []uint64) error {
	return db.DB.WithContext(ctx).Where("group_id = ? and user_id in ?", groupID, userIDs).Delete(&domain.GroupMember{}).Error
}

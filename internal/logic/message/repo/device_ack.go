package repo

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gim/internal/logic/message/domain"
	"gim/pkg/db"
)

type deviceACKRepo struct{}

var DeviceACKRepo = new(deviceACKRepo)

// Set 设置设备同步序列号
func (c *deviceACKRepo) Set(ctx context.Context, userID, deviceID, ack uint64) error {
	deviceACK := &domain.DeviceACK{
		DeviceID: deviceID,
		UserID:   userID,
		ACK:      ack,
	}
	return db.DB.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"ack", "updated_at"}),
	}).Create(deviceACK).Error
}

func (c *deviceACKRepo) List(ctx context.Context, userID uint64) ([]domain.DeviceACK, error) {
	return gorm.G[domain.DeviceACK](db.DB).Where("user_id = ?", userID).Find(ctx)
}

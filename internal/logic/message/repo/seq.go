package repo

import (
	"context"

	"gorm.io/gorm"

	"gim/pkg/db"
)

const (
	SeqObjectTypeUser = 1 // 用户
)

type seqRepo struct{}

var SeqRepo = new(seqRepo)

// Incr 自增seq,并且获取自增后的值
func (*seqRepo) Incr(ctx context.Context, objectType int, objectId uint64) (uint64, error) {
	var seq uint64
	err := db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// UPSERT: 插入或更新
		err := tx.Exec(
			"INSERT INTO seq (object_type, object_id, seq) VALUES (?, ?, 1) "+
				"ON DUPLICATE KEY UPDATE seq = seq + 1",
			objectType, objectId).Error
		if err != nil {
			return err
		}
		// 在同一事务中查询当前值
		return tx.Raw("SELECT seq FROM seq WHERE object_type = ? AND object_id = ?",
			objectType, objectId).Row().Scan(&seq)
	})
	if err != nil {
		return 0, err
	}
	return seq, nil
}

package repo

import (
	"context"

	"gim/pkg/db"
)

type seqRepo struct{}

var SeqRepo = new(seqRepo)

// Incr 自增seq,并且获取自增后的值
func (*seqRepo) Incr(ctx context.Context, userID uint64) (uint64, error) {
	const query = `
INSERT INTO seq (user_id, created_at, updated_at, seq)
VALUES (?, NOW(), NOW(), LAST_INSERT_ID(1))
ON DUPLICATE KEY UPDATE
  seq = LAST_INSERT_ID(seq + 1),
  updated_at = NOW()
`

	sqlDB, err := db.DB.DB()
	if err != nil {
		return 0, err
	}
	res, err := sqlDB.ExecContext(ctx, query, userID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

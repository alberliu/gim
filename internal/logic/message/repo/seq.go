package repo

import (
	"context"
	"database/sql"
	"errors"

	"gim/pkg/db"
)

const (
	SeqObjectTypeUser = 1 // 用户
	SeqObjectTypeRoom = 2 // 房间
)

type seqRepo struct{}

var SeqRepo = new(seqRepo)

// Incr 自增seq,并且获取自增后的值
func (*seqRepo) Incr(ctx context.Context, objectType int, objectId uint64) (uint64, error) {
	tx := db.DB.Begin()
	defer tx.Rollback()

	var seq uint64
	err := tx.WithContext(ctx).Raw("select seq from seq where object_type = ? and object_id = ? for update", objectType, objectId).
		Row().Scan(&seq)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		err = tx.WithContext(ctx).Exec("insert into seq (object_type,object_id,seq) values (?,?,?)", objectType, objectId, seq+1).Error
		if err != nil {
			return 0, err
		}
	} else {
		err = tx.WithContext(ctx).Exec("update seq set seq = seq + 1 where object_type = ? and object_id = ?", objectType, objectId).Error
		if err != nil {
			return 0, err
		}
	}

	tx.Commit()
	return seq + 1, nil
}

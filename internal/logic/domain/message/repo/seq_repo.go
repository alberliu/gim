package repo

import (
	"database/sql"
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

const (
	SeqObjectTypeUser = 1 // 用户
	SeqObjectTypeRoom = 2 // 房间
)

type seqRepo struct{}

var SeqRepo = new(seqRepo)

// Incr 自增seq,并且获取自增后的值
func (*seqRepo) Incr(objectType int, objectId int64) (int64, error) {
	tx := db.DB.Begin()
	defer tx.Rollback()

	var seq int64
	err := tx.Raw("select seq from seq where object_type = ? and object_id = ? for update", objectType, objectId).
		Row().Scan(&seq)
	if err != nil && err != sql.ErrNoRows {
		return 0, gerrors.WrapError(err)
	}
	if err == sql.ErrNoRows {
		err = tx.Exec("insert into seq (object_type,object_id,seq) values (?,?,?)", objectType, objectId, seq+1).Error
		if err != nil {
			return 0, gerrors.WrapError(err)
		}
	} else {
		err = tx.Exec("update seq set seq = seq + 1 where object_type = ? and object_id = ?", objectType, objectId).Error
		if err != nil {
			return 0, gerrors.WrapError(err)
		}
	}

	tx.Commit()
	return seq + 1, nil
}

package room

import (
	"fmt"
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

const RoomSeqKey = "room_seq:%d"

type roomSeqRepo struct{}

var RoomSeqRepo = new(roomSeqRepo)

func (*roomSeqRepo) GetNextSeq(roomId int64) (int64, error) {
	num, err := db.RedisCli.Incr(fmt.Sprintf(RoomSeqKey, roomId)).Result()
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return num, err
}

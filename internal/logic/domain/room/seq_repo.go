package room

import (
	"fmt"

	"gim/pkg/db"
	"gim/pkg/gerrors"
)

const SeqKey = "room_seq:%d"

type seqRepo struct{}

var SeqRepo = new(seqRepo)

func (*seqRepo) GetNextSeq(roomId int64) (int64, error) {
	num, err := db.RedisCli.Incr(fmt.Sprintf(SeqKey, roomId)).Result()
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return num, err
}

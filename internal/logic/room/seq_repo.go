package room

import (
	"fmt"

	"gim/pkg/db"
)

const SeqKey = "room_seq:%d"

type seqRepo struct{}

var SeqRepo = new(seqRepo)

func (*seqRepo) GetNextSeq(roomId uint64) (uint64, error) {
	num, err := db.RedisCli.Incr(fmt.Sprintf(SeqKey, roomId)).Result()
	if err != nil {
		return 0, err
	}
	return uint64(num), err
}

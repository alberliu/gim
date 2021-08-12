package cache

import (
	"fmt"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"gim/pkg/pb"
	"gim/pkg/util"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
)

const RoomMessageKey = "room_message:%d"

const RoomMessageExpireTime = 2 * time.Minute

type roomMessageCache struct{}

var RoomMessageCache = new(roomMessageCache)

// Add 将消息添加到队列
func (*roomMessageCache) Add(roomId int64, msg *pb.Message) error {
	key := fmt.Sprintf(RoomMessageKey, roomId)
	buf, err := proto.Marshal(msg)
	if err != nil {
		return gerrors.WrapError(err)
	}
	_, err = db.RedisCli.ZAdd(key, redis.Z{
		Score:  float64(msg.Seq),
		Member: buf,
	}).Result()

	db.RedisCli.Expire(key, RoomMessageExpireTime)
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// List 获取指定房间序列号大于seq的消息
func (*roomMessageCache) List(roomId int64, seq int64) ([]*pb.Message, error) {
	key := fmt.Sprintf(RoomMessageKey, roomId)
	result, err := db.RedisCli.ZRangeByScore(key, redis.ZRangeBy{
		Min: strconv.FormatInt(seq, 10),
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	var msgs []*pb.Message
	for i := range result {
		buf := util.Str2bytes(result[i])
		var msg pb.Message
		err = proto.Unmarshal(buf, &msg)
		if err != nil {
			return nil, gerrors.WrapError(err)
		}
		msgs = append(msgs, &msg)
	}
	return msgs, nil
}

func (*roomMessageCache) ListByIndex(roomId int64, start, stop int64) ([]*pb.Message, error) {
	key := fmt.Sprintf(RoomMessageKey, roomId)
	result, err := db.RedisCli.ZRange(key, start, stop).Result()
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	var msgs []*pb.Message
	for i := range result {
		buf := util.Str2bytes(result[i])
		var msg pb.Message
		err = proto.Unmarshal(buf, &msg)
		if err != nil {
			return nil, gerrors.WrapError(err)
		}
		msgs = append(msgs, &msg)
	}
	return msgs, nil
}

func (*roomMessageCache) DelBySeq(roomId int64, min, max int64) error {
	if min == 0 && max == 0 {
		return nil
	}
	key := fmt.Sprintf(RoomMessageKey, roomId)
	_, err := db.RedisCli.ZRemRangeByScore(key, strconv.FormatInt(min, 10), strconv.FormatInt(max, 10)).Result()
	return gerrors.WrapError(err)
}

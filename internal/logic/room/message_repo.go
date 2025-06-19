package room

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"google.golang.org/protobuf/proto"

	"gim/pkg/db"
	pb "gim/pkg/protocol/pb/logicpb"
	"gim/pkg/util"
)

const MessageKey = "room_message:%d"

const MessageExpireTime = 2 * time.Minute

type messageRepo struct{}

var MessageRepo = new(messageRepo)

// Add 将消息添加到队列
func (*messageRepo) Add(roomID uint64, msg *pb.Message) error {
	key := fmt.Sprintf(MessageKey, roomID)
	buf, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = db.RedisCli.ZAdd(key, redis.Z{
		Score:  float64(msg.Seq),
		Member: buf,
	}).Result()
	if err != nil {
		return err
	}

	_, err = db.RedisCli.Expire(key, MessageExpireTime).Result()
	return err
}

// List 获取指定房间序列号大于seq的消息
func (*messageRepo) List(roomID, seq uint64) ([]*pb.Message, error) {
	key := fmt.Sprintf(MessageKey, roomID)
	result, err := db.RedisCli.ZRangeByScore(key, redis.ZRangeBy{
		Min: strconv.FormatUint(seq, 10),
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil, err
	}

	var msgs []*pb.Message
	for i := range result {
		buf := util.Str2bytes(result[i])
		var msg pb.Message
		err = proto.Unmarshal(buf, &msg)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, &msg)
	}
	return msgs, nil
}

func (*messageRepo) ListByIndex(roomID uint64, start, stop int64) ([]*pb.Message, error) {
	key := fmt.Sprintf(MessageKey, roomID)
	result, err := db.RedisCli.ZRange(key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	var msgs []*pb.Message
	for i := range result {
		buf := util.Str2bytes(result[i])
		var msg pb.Message
		err = proto.Unmarshal(buf, &msg)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, &msg)
	}
	return msgs, nil
}

func (*messageRepo) DelBySeq(roomID uint64, min, max uint64) error {
	if min == 0 && max == 0 {
		return nil
	}
	key := fmt.Sprintf(MessageKey, roomID)
	_, err := db.RedisCli.ZRemRangeByScore(key, strconv.FormatUint(min, 10), strconv.FormatUint(max, 10)).Result()
	return err
}

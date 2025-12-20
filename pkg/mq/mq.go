package mq

import (
	"context"

	"gim/pkg/db"
)

const (
	PushRoomTopic         = "pushRoom"         // 房间消息队列
	PushRoomPriorityTopic = "pushRoomPriority" // 房间优先级消息队列
	PushAllTopic          = "pushAll"          // 全服消息队列
)

func Publish(ctx context.Context, topic string, bytes []byte) error {
	_, err := db.RedisCli.Publish(ctx, topic, bytes).Result()
	return err
}

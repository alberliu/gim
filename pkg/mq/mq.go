package mq

import (
	"context"

	"gim/pkg/db"
)

const (
	PushRoomTopic         = "push_room_topic"          // 房间消息队列
	PushRoomPriorityTopic = "push_room_priority_topic" // 房间优先级消息队列
	PushAllTopic          = "push_all_topic"           // 全服消息队列
)

func Publish(ctx context.Context, topic string, bytes []byte) error {
	_, err := db.RedisCli.Publish(ctx, topic, bytes).Result()
	return err
}

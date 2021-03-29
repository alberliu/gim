package connect

import (
	"gim/config"
	"gim/pkg/db"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/topic"

	"github.com/go-redis/redis"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

func StartSubscribe() {
	channel := db.RedisCli.Subscribe(topic.PushRoomTopic, topic.PushAllTopic).Channel()
	for i := 0; i < config.Connect.SubscribeNum; i++ {
		go handleMsg(channel)
	}
}

func handleMsg(channel <-chan *redis.Message) {
	for msg := range channel {
		if msg.Channel == topic.PushRoomTopic {
			handlePushRoom([]byte(msg.Payload))
		}
		if msg.Channel == topic.PushAllTopic {
			handlePushAll([]byte(msg.Payload))
		}
	}
}

func handlePushRoom(bytes []byte) {
	var msg pb.PushRoomMsg
	err := proto.Unmarshal(bytes, &msg)
	if err != nil {
		logger.Logger.Error("handlePushRoom error", zap.Error(err))
		return
	}
	PushRoom(msg.RoomId, msg.MessageSend)
}

func handlePushAll(bytes []byte) {
	var msg pb.PushAllMsg
	err := proto.Unmarshal(bytes, &msg)
	if err != nil {
		logger.Logger.Error("handlePushRoom error", zap.Error(err))
		return
	}
	PushAll(msg.MessageSend)
}

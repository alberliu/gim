package connect

import (
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"gim/config"
	"gim/pkg/db"
	"gim/pkg/logger"
	"gim/pkg/mq"
	"gim/pkg/protocol/pb"
)

// StartSubscribe 启动MQ消息处理逻辑
func StartSubscribe() {
	pushRoomPriorityChannel := db.RedisCli.Subscribe(mq.PushRoomPriorityTopic).Channel()
	pushRoomChannel := db.RedisCli.Subscribe(mq.PushRoomTopic).Channel()
	for i := 0; i < config.Config.PushRoomSubscribeNum; i++ {
		go handlePushRoomMsg(pushRoomPriorityChannel, pushRoomChannel)
	}

	pushAllChannel := db.RedisCli.Subscribe(mq.PushAllTopic).Channel()
	for i := 0; i < config.Config.PushAllSubscribeNum; i++ {
		go handlePushAllMsg(pushAllChannel)
	}
}

func handlePushRoomMsg(priorityChannel, channel <-chan *redis.Message) {
	for {
		select {
		case msg := <-priorityChannel:
			handlePushRoom([]byte(msg.Payload))
		default:
			select {
			case msg := <-channel:
				handlePushRoom([]byte(msg.Payload))
			default:
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}
	}
}

func handlePushAllMsg(channel <-chan *redis.Message) {
	for msg := range channel {
		handlePushAll([]byte(msg.Payload))
	}
}

func handlePushRoom(bytes []byte) {
	var msg pb.PushRoomMsg
	err := proto.Unmarshal(bytes, &msg)
	if err != nil {
		logger.Logger.Error("handlePushRoom error", zap.Error(err))
		return
	}
	PushRoom(msg.RoomId, msg.Message)
}

func handlePushAll(bytes []byte) {
	var msg pb.PushAllMsg
	err := proto.Unmarshal(bytes, &msg)
	if err != nil {
		logger.Logger.Error("handlePushRoom error", zap.Error(err))
		return
	}
	PushAll(msg.Message)
}

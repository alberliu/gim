package connect

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"

	"gim/config"
	"gim/pkg/db"
	"gim/pkg/mq"
	pb "gim/pkg/protocol/pb/connectpb"
)

// StartSubscribe 启动MQ消息处理逻辑
func StartSubscribe() {
	pushRoomPriorityChannel := db.RedisCli.Subscribe(context.Background(), mq.PushRoomPriorityTopic).Channel()
	pushRoomChannel := db.RedisCli.Subscribe(context.Background(), mq.PushRoomTopic).Channel()
	for i := 0; i < config.Config.PushRoomSubscribeNum; i++ {
		go handlePushRoomMsg(pushRoomPriorityChannel, pushRoomChannel)
	}

	pushAllChannel := db.RedisCli.Subscribe(context.Background(), mq.PushAllTopic).Channel()
	for i := 0; i < config.Config.PushAllSubscribeNum; i++ {
		go handlePushAllMsg(pushAllChannel)
	}
}

func handlePushRoomMsg(priorityChannel, channel <-chan *redis.Message) {
	for {
		var msg *redis.Message
		select {
		case msg = <-priorityChannel:
		default:
			select {
			case msg = <-priorityChannel:
			case msg = <-channel:
			}
		}
		handlePushRoom([]byte(msg.Payload))
	}
}

func handlePushAllMsg(channel <-chan *redis.Message) {
	for msg := range channel {
		handlePushAll([]byte(msg.Payload))
	}
}

func handlePushRoom(buf []byte) {
	var message pb.PushRoomMessage
	err := proto.Unmarshal(buf, &message)
	if err != nil {
		slog.Error("handlePushRoom error", "error", err)
		return
	}
	slog.Debug("handlePushRoom", "msg", &message)
	PushRoom(message.RoomId, message.Message)
}

func handlePushAll(buf []byte) {
	var msg pb.PushAllMessage
	err := proto.Unmarshal(buf, &msg)
	if err != nil {
		slog.Error("handlePushRoom error", "error", err)
		return
	}
	slog.Debug("handlePushAll", "msg", &msg)
	PushAll(msg.Message)
}

package connect

import (
	"log/slog"
	"time"

	"github.com/go-redis/redis"
	"google.golang.org/protobuf/proto"

	"gim/config"
	"gim/pkg/db"
	"gim/pkg/mq"
	pb "gim/pkg/protocol/pb/connectpb"
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
		slog.Error("handlePushRoom error", "error", err)
		return
	}
	slog.Debug("handlePushRoom", "msg", &msg)
	PushRoom(msg.RoomId, msg.Message)
}

func handlePushAll(bytes []byte) {
	var msg pb.PushAllMsg
	err := proto.Unmarshal(bytes, &msg)
	if err != nil {
		slog.Error("handlePushRoom error", "error", err)
		return
	}
	slog.Debug("handlePushAll", "msg", &msg)
	PushAll(msg.Message)
}

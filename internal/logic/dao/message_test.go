package dao

import (
	"fmt"
	"gim/internal/logic/model"
	"testing"
	"time"
)

func TestMessageDao_Add(t *testing.T) {
	message := model.Message{
		ObjectType:     1,
		ObjectId:       1,
		RequestId:      1,
		SenderType:     1,
		SenderId:       1,
		SenderDeviceId: 1,
		ReceiverType:   1,
		ReceiverId:     1,
		ToUserIds:      "1",
		Type:           1,
		Content:        []byte("123456"),
		Seq:            2,
		SendTime:       time.Now(),
		Status:         0,
	}
	fmt.Println(MessageDao.Add("message", message))
}

func TestMessageDao_ListByUserIdAndUserSeq(t *testing.T) {
	messages, err := MessageDao.ListBySeq("message", 1, 1, 0)
	fmt.Println(err)
	for i := range messages {
		fmt.Printf("%+v\n", messages[i])
	}
}

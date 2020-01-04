package dao

import (
	"encoding/json"
	"fmt"
	"gim/internal/logic/model"
	"testing"
	"time"
)

func TestMessageDao_Add(t *testing.T) {
	message := model.Message{
		AppId:          2,
		ObjectType:     1,
		ObjectId:       1,
		SenderType:     2,
		SenderId:       2,
		SenderDeviceId: 2,
		ReceiverType:   2,
		ReceiverId:     2,
		ToUserIds:      "2",
		Type:           1,
		Content:        "",
		Seq:            2,
		SendTime:       time.Now(),
	}
	fmt.Println(MessageDao.Add("message", message))
}

func TestMessageDao_ListByUserIdAndUserSeq(t *testing.T) {
	messages, err := MessageDao.ListBySeq("message", 1, 1, 1, 0)
	fmt.Println(err)
	bytes, _ := json.Marshal(messages)
	fmt.Println(string(bytes))
}

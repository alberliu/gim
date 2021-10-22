package repo

import (
	"fmt"
	"gim/internal/logic/domain/message/model"
	"testing"
	"time"
)

func TestMessageDao_Add(t *testing.T) {
	message := model.Message{
		UserId:       1,
		RequestId:    1,
		SenderType:   1,
		SenderId:     1,
		ReceiverType: 1,
		ReceiverId:   1,
		ToUserIds:    "1",
		Type:         1,
		Content:      []byte("123456"),
		Seq:          2,
		SendTime:     time.Now(),
		Status:       0,
	}
	fmt.Println(MessageRepo.Save(message))
}

func TestMessageDao_ListByUserIdAndUserSeq(t *testing.T) {
	messages, hasMore, err := MessageRepo.ListBySeq(1, 0, 100)
	fmt.Println(err)
	fmt.Println(hasMore)
	for i := range messages {
		fmt.Printf("%+v\n", messages[i])
	}
}

func Test_messageDao_tableName(t *testing.T) {
	fmt.Println(MessageRepo.tableName(1001))
}

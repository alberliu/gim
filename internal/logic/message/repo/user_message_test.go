package repo

import (
	"testing"

	"gim/internal/logic/message/domain"
)

func TestUserMessageDao_Add(t *testing.T) {
	message := domain.UserMessage{
		UserID:    1,
		Seq:       1,
		MessageID: 1,
	}
	t.Log(UserMessageRepo.Create(&message))
}

func TestUserMessageDao_ListByUserIdAndUserSeq(t *testing.T) {
	messages, hasMore, err := UserMessageRepo.ListBySeq(1, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hasMore)
	for i := range messages {
		t.Logf("%+v\n", messages[i])
	}
}

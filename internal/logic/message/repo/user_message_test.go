package repo

import (
	"context"
	"testing"

	"gim/internal/logic/message/domain"
)

func TestUserMessageDao_Add(t *testing.T) {
	message := domain.UserMessage{
		UserID:    1,
		Seq:       1,
		MessageID: 1,
	}
	t.Log(UserMessageRepo.Create(context.Background(), &message))
}

func TestUserMessageDao_ListByUserIdAndUserSeq(t *testing.T) {
	messages, hasMore, err := UserMessageRepo.ListBySeq(context.Background(), 1, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hasMore)
	for i := range messages {
		t.Logf("%+v\n", messages[i])
	}
}

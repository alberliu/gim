package repo

import (
	"context"
	"testing"
)

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

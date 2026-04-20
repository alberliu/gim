package repo

import (
	"context"
	"testing"

	"gim/internal/logic/message/domain"
)

func Test_messageRepo_Save(t *testing.T) {
	msg := domain.Message{
		RequestID: "1",
		Command:   1,
		Content:   []byte("hello world"),
	}
	err := MessageRepo.Save(context.Background(), &msg)
	t.Log(err)
}

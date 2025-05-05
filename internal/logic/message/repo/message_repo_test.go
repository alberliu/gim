package repo

import (
	"testing"

	"gim/internal/logic/message/domain"
)

func Test_messageRepo_Save(t *testing.T) {
	msg := domain.Message{
		RequestID: 1,
		Code:      1,
		Content:   []byte("hello world"),
		Status:    0,
	}
	err := MessageRepo.Save(&msg)
	t.Log(err)
}

func Test_messageRepo_GetByIDs(t *testing.T) {
	msgs, err := MessageRepo.GetByIDs([]int64{1})
	t.Log(err)
	t.Log(msgs)
}

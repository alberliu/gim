package service

import (
	"context"
	"fmt"
	"testing"
)

func TestUserRequenceService_GetNext(t *testing.T) {
}

func Test_seqService_GetRoomNext(t *testing.T) {
	seq, err := SeqService.GetRoomNext(context.TODO(), 1)
	fmt.Println(seq, err)
}

func Test_seqService_GetUserNext(t *testing.T) {
	seq, err := SeqService.GetUserNext(context.TODO(), 1)
	fmt.Println(seq, err)
}

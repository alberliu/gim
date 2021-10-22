package service

import (
	"context"
	"fmt"
	"testing"
)

func Test_seqService_GetUserNext(t *testing.T) {
	seq, err := SeqService.GetUserNext(context.TODO(), 1)
	fmt.Println(seq, err)
}

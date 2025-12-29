package repo

import (
	"context"
	"testing"
)

func Test_seqDao_Incr(t *testing.T) {
	seq, err := SeqRepo.Incr(context.Background(), 1, 5)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(seq)
}

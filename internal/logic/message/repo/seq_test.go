package repo

import (
	"context"
	"fmt"
	"testing"
)

func Test_seqDao_Incr(t *testing.T) {
	fmt.Println(SeqRepo.Incr(context.Background(), 1, 5))
}

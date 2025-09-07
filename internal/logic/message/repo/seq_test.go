package repo

import (
	"fmt"
	"testing"
)

func Test_seqDao_Incr(t *testing.T) {
	fmt.Println(SeqRepo.Incr(1, 5))
}

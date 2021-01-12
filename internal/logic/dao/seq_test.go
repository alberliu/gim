package dao

import (
	"fmt"
	"testing"
)

func Test_seqDao_Incr(t *testing.T) {
	fmt.Println(SeqDao.Incr(1, 5))
}

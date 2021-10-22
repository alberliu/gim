package repo

import (
	"fmt"
	"testing"
)

func TestGroupDao_Get(t *testing.T) {
	group, err := GroupDao.Get(1)
	fmt.Printf("%+v\n %+v\n", group, err)
}

package repo

import (
	"testing"
)

func TestGroupDao_Get(t *testing.T) {
	group, err := GroupDao.Get(5)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(group)

	for i := range group.Members {
		t.Log(group.Members[i])
	}
}

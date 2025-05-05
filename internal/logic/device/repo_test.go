package device

import (
	"testing"
)

func Test_repo_Get(t *testing.T) {
	device, err := Repo.Get(1)
	t.Log(err)
	t.Log(device)
}

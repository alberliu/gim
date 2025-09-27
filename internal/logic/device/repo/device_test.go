package repo

import (
	"context"
	"testing"
)

func Test_deviceRepo_Get(t *testing.T) {
	device, err := DeviceRepo.Get(context.Background(), 1)
	t.Log(err)
	t.Log(device)
}

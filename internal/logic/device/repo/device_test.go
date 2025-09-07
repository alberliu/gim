package repo

import (
	"testing"
)

func Test_deviceRepo_Get(t *testing.T) {
	device, err := DeviceRepo.Get(1)
	t.Log(err)
	t.Log(device)
}

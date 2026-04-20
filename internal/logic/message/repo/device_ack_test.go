package repo

import (
	"context"
	"testing"
)

func Test_deviceACKRepo_Set(t *testing.T) {
	err := DeviceACKRepo.Set(context.TODO(), 1, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_deviceACKRepo_List(t *testing.T) {
	acks, err := DeviceACKRepo.List(context.TODO(), 1)
	if err != nil {
		t.Fatal(err)
	}
	for _, ack := range acks {
		t.Log(ack)
	}
}

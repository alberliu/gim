package dao

import (
	"fmt"
	"testing"
)

func TestDeviceAckDao_Add(t *testing.T) {
	fmt.Println(DeviceAckDao.Add(ctx, 1, 1))
}

func TestDeviceAckDao_Get(t *testing.T) {
	fmt.Println(DeviceAckDao.Get(ctx, 1))
}

func TestDeviceAckDao_Update(t *testing.T) {
	fmt.Println(DeviceAckDao.Update(ctx, 1, 2))
}

func TestDeviceAckDao_GetMaxByUserId(t *testing.T) {
	fmt.Println(DeviceAckDao.GetMaxByUserId(ctx, 1, 2))
}

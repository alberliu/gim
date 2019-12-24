package dao

import (
	"fmt"
	"testing"
)

func TestDeviceAckDao_Add(t *testing.T) {
	fmt.Println(DeviceAckDao.Add(1, 1))
}

func TestDeviceAckDao_Get(t *testing.T) {
	fmt.Println(DeviceAckDao.Get(1))
}

func TestDeviceAckDao_Update(t *testing.T) {
	fmt.Println(DeviceAckDao.Update(1, 2))
}

func TestDeviceAckDao_GetMaxByUserId(t *testing.T) {
	fmt.Println(DeviceAckDao.GetMaxByUserId(1, 2))
}

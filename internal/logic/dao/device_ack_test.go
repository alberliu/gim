package dao

import (
	"fmt"
	"testing"
)

func TestDeviceAckDao_Add(t *testing.T) {
	fmt.Println(DeviceAckDao.Add(10, 10))
}

func TestDeviceAckDao_Get(t *testing.T) {
	fmt.Println(DeviceAckDao.Get(10))
}

func TestDeviceAckDao_Update(t *testing.T) {
	fmt.Println(DeviceAckDao.Update(10, 12))
}

func TestDeviceAckDao_GetMaxByUserId(t *testing.T) {
	fmt.Println(DeviceAckDao.GetMaxByUserId(1))
}

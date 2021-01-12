package dao

import (
	"fmt"
	"gim/internal/logic/model"
	"testing"
)

func init() {
	fmt.Println("start")
}

func TestDeviceDao_Add(t *testing.T) {
	device := model.Device{
		UserId:        1,
		Type:          1,
		Brand:         "huawei",
		Model:         "huawei P10",
		SystemVersion: "8.0.0",
		SDKVersion:    "1.0.0",
		Status:        1,
	}
	fmt.Println(DeviceDao.Add(device))
}

func TestDeviceDao_Get(t *testing.T) {
	device, err := DeviceDao.Get(1)
	fmt.Printf("%+v\n %+v\n", device, err)
}

func TestDeviceDao_ListOnlineByUserId(t *testing.T) {
	devices, err := DeviceDao.ListOnlineByUserId(1)
	fmt.Println(err)
	fmt.Printf("%+v \n", devices)
}

func TestDeviceDao_UpdateStatus(t *testing.T) {
	fmt.Println(DeviceDao.UpdateStatus(1, 0))
}

func TestDeviceDao_Upgrade(t *testing.T) {
	fmt.Println(DeviceDao.Upgrade(1, "9.0.0", "2.0.0"))
}

package device

import (
	"fmt"
	"testing"
)

func TestDeviceDao_Add(t *testing.T) {
	device := Device{
		UserId:        1,
		Type:          1,
		Brand:         "huawei",
		Model:         "huawei P10",
		SystemVersion: "8.0.0",
		SDKVersion:    "1.0.0",
		Status:        1,
	}
	err := Dao.Save(&device)
	fmt.Println(err)
	fmt.Println(device)
}

func TestDeviceDao_Get(t *testing.T) {
	device, err := Dao.Get(1)
	fmt.Printf("%+v\n %+v\n", device, err)
}

func TestDeviceDao_ListOnlineByUserId(t *testing.T) {
	devices, err := Dao.ListOnlineByUserId(1)
	fmt.Println(err)
	fmt.Printf("%+v \n", devices)
}

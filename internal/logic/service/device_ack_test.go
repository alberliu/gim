package service

import (
	"context"
	"fmt"
	_ "gim/internal/logic/dao"
	"gim/pkg/db"
	"testing"
)

func Test_deviceAckService_GetMaxByUserId(t *testing.T) {
	db.InitByTest()
	fmt.Println(DeviceAckService.Update(context.TODO(), 1, 2, 2))
}

func Test_deviceAckService_Update(t *testing.T) {
	db.InitByTest()
	fmt.Println(DeviceAckService.GetMaxByUserId(context.TODO(), 1))
}

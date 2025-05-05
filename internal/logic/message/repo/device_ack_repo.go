package repo

import (
	"fmt"
	"strconv"

	"gim/pkg/db"
)

const DeviceACKKey = "device_ack:%d"

type deviceACKRepo struct{}

var DeviceACKRepo = new(deviceACKRepo)

// Set 设置设备同步序列号
func (c *deviceACKRepo) Set(userId, deviceId, ack uint64) error {
	key := fmt.Sprintf(DeviceACKKey, userId)
	_, err := db.RedisCli.HSet(key, strconv.FormatUint(deviceId, 10), strconv.FormatUint(ack, 10)).Result()
	return err
}

func (c *deviceACKRepo) Get(userId uint64) (map[uint64]uint64, error) {
	key := fmt.Sprintf(DeviceACKKey, userId)
	result, err := db.RedisCli.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}

	acks := make(map[uint64]uint64, len(result))
	for k, v := range result {
		deviceId, _ := strconv.ParseUint(k, 10, 64)
		ack, _ := strconv.ParseUint(v, 10, 64)
		acks[deviceId] = ack
	}
	return acks, nil
}

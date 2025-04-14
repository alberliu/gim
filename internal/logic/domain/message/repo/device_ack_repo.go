package repo

import (
	"strconv"

	"gim/pkg/db"
	"gim/pkg/gerrors"
)

const (
	DeviceACKKey = "device_ack:"
)

type deviceACKRepo struct{}

var DeviceACKRepo = new(deviceACKRepo)

// Set 设置设备同步序列号
func (c *deviceACKRepo) Set(userId int64, deviceId int64, ack int64) error {
	_, err := db.RedisCli.HSet(DeviceACKKey+strconv.FormatInt(userId, 10), strconv.FormatInt(deviceId, 10),
		strconv.FormatInt(ack, 10)).Result()

	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

func (c *deviceACKRepo) Get(userId int64) (map[int64]int64, error) {
	result, err := db.RedisCli.HGetAll(DeviceACKKey + strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	acks := make(map[int64]int64, len(result))
	for k, v := range result {
		deviceId, _ := strconv.ParseInt(k, 10, 64)
		ack, _ := strconv.ParseInt(v, 10, 64)
		acks[deviceId] = ack
	}
	return acks, nil
}

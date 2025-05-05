package repo

import (
	"encoding/json"
	"fmt"
	"strconv"

	"gim/internal/user/domain"
	"gim/pkg/db"
	"gim/pkg/util"
)

const AuthKey = "auth:%d"

type authRepo struct{}

var AuthRepo = new(authRepo)

func (*authRepo) Get(userId, deviceId uint64) (*domain.Device, error) {
	key := fmt.Sprintf(AuthKey, userId)
	bytes, err := db.RedisCli.HGet(key, strconv.FormatUint(deviceId, 10)).Bytes()
	if err != nil {
		return nil, err
	}

	var device domain.Device
	err = json.Unmarshal(bytes, &device)
	return &device, err
}

func (*authRepo) Set(userId, deviceId uint64, device domain.Device) error {
	bytes, err := json.Marshal(device)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(AuthKey, userId)
	_, err = db.RedisCli.HSet(key, strconv.FormatUint(deviceId, 10), bytes).Result()
	return err
}

func (*authRepo) GetAll(userId uint64) (map[uint64]domain.Device, error) {
	key := fmt.Sprintf(AuthKey, userId)
	result, err := db.RedisCli.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}

	var devices = make(map[uint64]domain.Device, len(result))

	for k, v := range result {
		deviceId, err := strconv.ParseUint(k, 10, 64)
		if err != nil {
			return nil, err
		}

		var device domain.Device
		err = json.Unmarshal(util.Str2bytes(v), &device)
		if err != nil {
			return nil, err
		}
		devices[deviceId] = device
	}
	return devices, nil
}

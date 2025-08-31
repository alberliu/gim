package repo

import (
	"encoding/json"
	"fmt"
	"strconv"

	"gim/internal/business/user/domain"
	"gim/pkg/db"
	"gim/pkg/util"
)

const AuthKey = "auth:%d"

type authRepo struct{}

var AuthRepo = new(authRepo)

func (*authRepo) Get(userID, deviceID uint64) (*domain.Device, error) {
	key := fmt.Sprintf(AuthKey, userID)
	bytes, err := db.RedisCli.HGet(key, strconv.FormatUint(deviceID, 10)).Bytes()
	if err != nil {
		return nil, err
	}

	var device domain.Device
	err = json.Unmarshal(bytes, &device)
	return &device, err
}

func (*authRepo) Set(userID, deviceID uint64, device domain.Device) error {
	bytes, err := json.Marshal(device)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(AuthKey, userID)
	_, err = db.RedisCli.HSet(key, strconv.FormatUint(deviceID, 10), bytes).Result()
	return err
}

func (*authRepo) GetAll(userID uint64) (map[uint64]domain.Device, error) {
	key := fmt.Sprintf(AuthKey, userID)
	result, err := db.RedisCli.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}

	var devices = make(map[uint64]domain.Device, len(result))

	for k, v := range result {
		deviceID, err := strconv.ParseUint(k, 10, 64)
		if err != nil {
			return nil, err
		}

		var device domain.Device
		err = json.Unmarshal(util.Str2bytes(v), &device)
		if err != nil {
			return nil, err
		}
		devices[deviceID] = device
	}
	return devices, nil
}

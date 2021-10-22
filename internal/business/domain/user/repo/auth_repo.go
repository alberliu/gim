package repo

import "gim/internal/business/domain/user/model"

type authRepo struct{}

var AuthRepo = new(authRepo)

func (*authRepo) Get(userId, deviceId int64) (*model.Device, error) {
	return AuthCache.Get(userId, deviceId)
}

func (*authRepo) Set(userId, deviceId int64, device model.Device) error {
	return AuthCache.Set(userId, deviceId, device)
}

func (*authRepo) GetAll(userId int64) (map[int64]model.Device, error) {
	return AuthCache.GetAll(userId)
}

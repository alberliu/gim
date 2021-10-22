package device

type deviceRepo struct{}

var DeviceRepo = new(deviceRepo)

// Get 获取设备
func (*deviceRepo) Get(deviceId int64) (*Device, error) {
	device, err := DeviceDao.Get(deviceId)
	if err != nil {
		return nil, err
	}
	return device, nil
}

// Save 保存设备信息
func (*deviceRepo) Save(device *Device) error {
	err := DeviceDao.Save(device)
	if err != nil {
		return err
	}

	if device.UserId != 0 {
		err = UserDeviceCache.Del(device.UserId)
		if err != nil {
			return err
		}
	}
	return nil
}

// ListOnlineByUserId 获取用户的所有在线设备
func (*deviceRepo) ListOnlineByUserId(userId int64) ([]Device, error) {
	devices, err := UserDeviceCache.Get(userId)
	if err != nil {
		return nil, err
	}

	if devices != nil {
		return devices, nil
	}

	devices, err = DeviceDao.ListOnlineByUserId(userId)
	if err != nil {
		return nil, err
	}

	err = UserDeviceCache.Set(userId, devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// ListOnlineByConnAddr 查询用户所有的在线设备
func (*deviceRepo) ListOnlineByConnAddr(connAddr string) ([]Device, error) {
	return DeviceDao.ListOnlineByConnAddr(connAddr)
}

// UpdateStatusOffline 更新设备为离线状态
func (*deviceRepo) UpdateStatusOffline(device Device) error {
	affected, err := DeviceDao.UpdateStatus(device.Id, device.ConnAddr, DeviceOffLine)
	if err != nil {
		return err
	}

	if affected == 1 && device.UserId != 0 {
		err = UserDeviceCache.Del(device.UserId)
		if err != nil {
			return err
		}
	}

	return nil
}

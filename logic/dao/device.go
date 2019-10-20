package dao

import (
	"database/sql"
	"gim/logic/db"
	"gim/logic/model"
	"gim/public/imctx"
	"gim/public/logger"
)

type deviceDao struct{}

var DeviceDao = new(deviceDao)

// Insert 插入一条设备信息
func (*deviceDao) Add(ctx *imctx.Context, device model.Device) error {
	_, err := db.DBCli.Exec(`insert into device(device_id,app_id,type,brand,model,system_version,sdk_version,status) 
		values(?,?,?,?,?,?,?,?)`,
		device.DeviceId, device.AppId, device.Type, device.Brand, device.Model, device.SystemVersion, device.SDKVersion, device.Status)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// Get 获取设备
func (*deviceDao) Get(ctx *imctx.Context, deviceId int64) (*model.Device, error) {
	device := model.Device{
		DeviceId: deviceId,
	}
	row := db.DBCli.QueryRow(`
		select app_id,user_id,type,brand,model,system_version,sdk_version,status,create_time,update_time
		from device where device_id = ?`, deviceId)
	err := row.Scan(&device.AppId, &device.UserId, &device.Type, &device.Brand, &device.Model, &device.SystemVersion, &device.SDKVersion,
		&device.Status, &device.CreateTime, &device.UpdateTime)
	if err != nil && err != sql.ErrNoRows {
		logger.Sugar.Error(err)
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &device, err
}

// ListUserOnline 查询用户所有的在线设备
func (*deviceDao) ListOnlineByUserId(ctx *imctx.Context, appId, userId int64) ([]model.Device, error) {
	rows, err := db.DBCli.Query(
		`select device_id,type,brand,model,system_version,sdk_version,status,create_time,update_time from device where app_id = ? and user_id = ? and status = ?`,
		appId, userId, model.DeviceOnLine)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	devices := make([]model.Device, 0, 5)
	for rows.Next() {
		device := new(model.Device)
		err = rows.Scan(&device.DeviceId, &device.Type, &device.Brand, &device.Model, &device.SystemVersion, &device.SDKVersion,
			&device.Status, &device.CreateTime, &device.UpdateTime)
		if err != nil {
			logger.Sugar.Error(err)
			return nil, err
		}
		devices = append(devices, *device)
	}
	return devices, nil
}

// UpdateUserIdAndStatus 更新设备绑定用户和设备在线状态
func (*deviceDao) UpdateUserIdAndStatus(ctx *imctx.Context, deviceId, userId int64, status int) error {
	_, err := db.DBCli.Exec("update device  set user_id = ?,status = ? where  device_id = ? ",
		userId, status, deviceId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// UpdateStatus 更新设备的在线状态
func (*deviceDao) UpdateStatus(ctx *imctx.Context, deviceId int64, status int) error {
	_, err := db.DBCli.Exec("update device set status = ? where device_id = ?", status, deviceId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// Upgrade 升级设备
func (*deviceDao) Upgrade(ctx *imctx.Context, deviceId int64, systemVersion, sdkVersion string) error {
	_, err := db.DBCli.Exec("update device set system_version = ?,sdk_version = ? where device_id = ? ",
		systemVersion, sdkVersion, deviceId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

package dao

import (
	"gim/internal/logic/model"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"time"

	"github.com/jinzhu/gorm"
)

type deviceDao struct{}

var DeviceDao = new(deviceDao)

// Add 插入一条设备信息
func (*deviceDao) Add(device model.Device) (int64, error) {
	device.CreateTime = time.Now()
	device.UpdateTime = time.Now()
	err := db.DB.Create(&device).Error
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return device.Id, nil
}

// Get 获取设备
func (*deviceDao) Get(deviceId int64) (*model.Device, error) {
	var device = model.Device{Id: deviceId}
	err := db.DB.First(&device).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, gerrors.WrapError(err)
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &device, nil
}

// ListOnlineByUserId 查询用户所有的在线设备
func (*deviceDao) ListOnlineByUserId(userId int64) ([]model.Device, error) {
	var devices []model.Device
	err := db.DB.Find(&devices, "user_id = ? and status = ?", userId, model.DeviceOnLine).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return devices, nil
}

// Update 更新设备绑定用户和设备在线状态
func (*deviceDao) Update(deviceId, userId int64, status int, connAddr string, clientAddr string) error {
	err := db.DB.Exec("update device set user_id = ?,status = ?,conn_addr = ?,client_addr = ? where id = ? ",
		userId, status, connAddr, clientAddr, deviceId).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// UpdateStatus 更新设备的在线状态
func (*deviceDao) UpdateStatus(deviceId int64, status int) error {
	err := db.DB.Exec("update device set status = ? where id = ?", status, deviceId).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Upgrade 升级设备
func (*deviceDao) Upgrade(deviceId int64, systemVersion, sdkVersion string) error {
	err := db.DB.Exec("update device set system_version = ?,sdk_version = ? where id = ? ",
		systemVersion, sdkVersion, deviceId).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// ListOnlineByConnAddr 查询用户所有的在线设备
func (*deviceDao) ListOnlineByConnAddr(connAddr string) ([]model.Device, error) {
	var devices []model.Device
	err := db.DB.Find(&devices, "conn_addr = ? and status = ?", connAddr, model.DeviceOnLine).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return devices, nil
}

// UpdateStatusByCoonAddr 更新在线状态
func (*deviceDao) UpdateStatusByCoonAddr(connAddr string, status int) error {
	err := db.DB.Model(&model.Device{}).Where("conn_addr = ?", connAddr).
		Update("status", status).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

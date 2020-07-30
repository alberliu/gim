package dao

import (
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

type deviceAckDao struct{}

var DeviceAckDao = new(deviceAckDao)

// Add 添加设备同步序列号记录
func (*deviceAckDao) Add(deviceId int64, ack int64) error {
	err := db.DB.Exec("insert into device_ack(device_id,ack) values(?,?)", deviceId, ack).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Get 获取设备同步序列号
func (*deviceAckDao) Get(deviceId int64) (int64, error) {
	row := db.DB.Raw("select ack from device_ack where device_id = ?", deviceId).Row()
	var ack int64
	err := row.Scan(&ack)
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return ack, nil
}

// UpdateSyncSequence 更新设备同步序列号
func (*deviceAckDao) Update(deviceId, ack int64) error {
	err := db.DB.Exec("update device_ack set ack = ? where device_id = ?", ack, deviceId).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// GetMaxByUserId 获取用户最大的同步序列号
func (*deviceAckDao) GetMaxByUserId(userId int64) (int64, error) {
	row := db.DB.Raw(`
		select max(a.ack) 
		from device d
		left join device_ack a on d.id = a.device_id  
		where d.user_id = ?`, userId).Row()
	var ack int64
	err := row.Scan(&ack)
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return ack, nil
}

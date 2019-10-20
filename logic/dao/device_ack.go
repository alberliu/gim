package dao

import (
	"gim/logic/db"
	"gim/public/imctx"
	"gim/public/logger"
)

type deviceAckDao struct{}

var DeviceAckDao = new(deviceAckDao)

// Add 添加设备同步序列号记录
func (*deviceAckDao) Add(ctx *imctx.Context, deviceId int64, ack int64) error {
	_, err := db.DBCli.Exec("insert into device_ack(device_id,ack) values(?,?)", deviceId, ack)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// Get 获取设备同步序列号
func (*deviceAckDao) Get(ctx *imctx.Context, deviceId int64) (int64, error) {
	row := db.DBCli.QueryRow("select ack from device_ack where device_id = ?", deviceId)
	var ack int64
	err := row.Scan(&ack)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}
	return ack, nil
}

// UpdateSyncSequence 更新设备同步序列号
func (*deviceAckDao) Update(ctx *imctx.Context, deviceId, ack int64) error {
	_, err := db.DBCli.Exec("update device_ack set ack = ? where device_id = ?", ack, deviceId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// GetMaxByUserId 获取用户最大的同步序列号
func (*deviceAckDao) GetMaxByUserId(ctx *imctx.Context, appId, userId int64) (int64, error) {
	row := db.DBCli.QueryRow(`
		select max(a.ack) 
		from device d
		left join device_ack a d.device_id = a.device_id  
		where d.app_id = ? and d.user_id = ?`, appId, userId)
	var ack int64
	err := row.Scan(&ack)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}
	return ack, nil
}

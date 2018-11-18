package dao

import (
	"goim/public/imctx"
	"goim/public/logger"
)

type deviceSyncSequenceDao struct{}

var DeviceSyncSequenceDao = new(deviceSyncSequenceDao)

// Add 添加设备同步序列号记录
func (*deviceSyncSequenceDao) Add(ctx *imctx.Context, deviceId int64, syncSequence int64) error {
	_, err := ctx.Session.Exec("insert into t_device_sync_sequence(device_id,sync_sequence) values(?,?)",
		deviceId, syncSequence)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// Get 获取设备同步序列号
func (*deviceSyncSequenceDao) Get(ctx *imctx.Context, id int64) (int64, error) {
	row := ctx.Session.QueryRow("select sync_sequence from t_device_sync_sequence where device_id = ?", id)
	var syncSequence int64
	err := row.Scan(&syncSequence)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}
	return syncSequence, nil
}

// GetMaxSyncSequenceByUserId 获取用户最大的同步序列号
func (*deviceSyncSequenceDao) GetMaxSyncSequenceByUserId(ctx *imctx.Context, userId int64) (int64, error) {
	row := ctx.Session.QueryRow(`
		select max(s.sync_sequence) 
		from t_device d
		left join t_device_sync_sequence s on d.id = s.device_id 
		where user_id = ?`, userId)
	var syncSequence int64
	err := row.Scan(&syncSequence)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}
	return syncSequence, nil
}

// UpdateSyncSequence 更新设备同步序列号
func (*deviceSyncSequenceDao) UpdateSyncSequence(ctx *imctx.Context, deviceId int64, sequence int64) error {
	_, err := ctx.Session.Exec("update t_device_sync_sequence set sync_sequence = ? where device_id = ?",
		sequence, deviceId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

package dao

import (
	"goim/public/imctx"
	"goim/public/logger"
)

type deviceSendSequenceDao struct{}

var DeviceSendSequenceDao = new(deviceSendSequenceDao)

// Add 添加设备发送序列号
func (*deviceSendSequenceDao) Add(ctx *imctx.Context, deviceId int64, sendSequence int64) error {
	_, err := ctx.Session.Exec("insert into t_device_send_sequence(device_id,send_sequence) values(?,?)",
		deviceId, sendSequence)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// Get 获取设备发送序列号
func (*deviceSendSequenceDao) Get(ctx *imctx.Context, id int64) (int64, error) {
	row := ctx.Session.QueryRow("select send_sequence from t_device_send_sequence where device_id = ?", id)
	var syncSeq int64
	err := row.Scan(&syncSeq)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, err
	}
	return syncSeq, nil
}

// UpdateSendSequence 更新设备发送序列号
func (*deviceSendSequenceDao) UpdateSendSequence(ctx *imctx.Context, deviceId int64, sequence int64) error {
	_, err := ctx.Session.Exec("update t_device_send_sequence set send_sequence = ? where device_id = ?",
		sequence, deviceId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

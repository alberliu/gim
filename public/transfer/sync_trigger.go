package transfer

// 同步消息触发
type SyncTrigger struct {
	DeviceId     int64 `json:"device_id"`     // 设备id
	UserId       int64 `json:"user_id"`       // 用户id
	SyncSequence int64 `json:"sync_sequence"` // 已经同步的消息序列号
}

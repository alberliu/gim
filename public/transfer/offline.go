package transfer

type OfflineReq struct {
	AppId    int64
	DeviceId int64 // 设备id
	UserId   int64 // 用户id
}

type OfflineResp struct {
}

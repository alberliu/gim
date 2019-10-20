package transfer

type MessageAckReq struct {
	IsSignIn bool // 标记用户是否登录成功
	AppId    int64
	DeviceId int64  // 设备id
	UserId   int64  // 用户id
	Bytes    []byte // 字节数组
}

type MessageAckResp struct {
}

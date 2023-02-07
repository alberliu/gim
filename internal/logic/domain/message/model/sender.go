package model

type Sender struct {
	UserId    int64  // 发送者id
	DeviceId  int64  // 发送者设备id
	Nickname  string // 昵称
	AvatarUrl string // 头像
	Extra     string // 扩展字段
}

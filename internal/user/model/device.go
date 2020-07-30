package model

type Device struct {
	Type   int32  // 设备类型,1:Android；2：IOS；3：Windows; 4：MacOS；5：Web
	Token  string // token
	Expire int64  // 过期时间
}

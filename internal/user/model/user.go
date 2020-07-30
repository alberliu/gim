package model

import (
	"time"
)

// User 账户
type User struct {
	Id          int64     // 用户id
	PhoneNumber string    // 手机号
	Nickname    string    // 昵称
	Sex         int32     // 性别，1:男；2:女
	AvatarUrl   string    // 用户头像
	Extra       string    // 附加属性
	CreateTime  time.Time // 创建时间
	UpdateTime  time.Time // 更新时间
}

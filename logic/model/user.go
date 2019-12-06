package model

import (
	"time"
)

// User 账户
type User struct {
	Id         int64     `json:"-"`          // 用户id
	AppId      int64     `json:"app_id"`     // app_id
	UserId     int64     `json:"user_id"`    // 手机号
	Nickname   string    `json:"nickname"`   // 昵称
	Sex        int32     `json:"sex"`        // 性别，1:男；2:女
	AvatarUrl  string    `json:"avatar_url"` // 用户头像
	Extra      string    `json:"extra"`      // 附加属性
	CreateTime time.Time `json:"-"`          // 创建时间
	UpdateTime time.Time `json:"-"`          // 更新时间
}

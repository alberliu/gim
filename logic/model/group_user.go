package model

import "time"

type GroupUser struct {
	Id         int64     `json:"id,omitempty"` // 自增主键
	AppId      int64     `json:"app_id"`       // app_id
	GroupId    int64     `json:"group_id"`     // 群组id
	UserId     int64     `json:"user_id"`      // 用户id
	Label      string    `json:"label"`        // 用户标签
	Extra      string    `json:"extra"`        // 群组用户附件属性
	CreateTime time.Time `json:"-"`            // 创建时间
	UpdateTime time.Time `json:"-"`            // 更新时间
}

// GroupUser 群组成员
type GroupUserInfo struct {
	UserId         int64     `json:"user_id"`          // 用户id
	Label          string    `json:"label"`            // 用户标签
	GroupUserExtra string    `json:"group_user_extra"` // 群组用户附件属性
	Nickname       string    `json:"name"`             // 昵称
	Sex            int       `json:"sex"`              // 性别,0:位置；1:男；2:女
	AvatarUrl      string    `json:"img"`              // 用户头像
	UserExtra      string    `json:"user_extra"`       // 用户附件属性
	CreateTime     time.Time `json:"create_time"`      // 创建时间
	UpdateTime     time.Time `json:"update_time"`      // 更新时间
}

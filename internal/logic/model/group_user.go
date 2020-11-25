package model

import "time"

type GroupUser struct {
	Id         int64     `json:"id,omitempty"` // 自增主键
	GroupId    int64     `json:"group_id"`     // 群组id
	UserId     int64     `json:"user_id"`      // 用户id
	MemberType int       `json:"member_type"`  // 群组类型
	Remarks    string    `json:"remarks"`      // 备注
	Extra      string    `json:"extra"`        // 附加属性
	Status     int       `json:"status"`       // 状态
	CreateTime time.Time `json:"-"`            // 创建时间
	UpdateTime time.Time `json:"-"`            // 更新时间
}

// GroupUser 群组成员
type GroupUserInfo struct {
	UserId     int64  `json:"user_id"`     // 用户id
	Nickname   string `json:"name"`        // 昵称
	Sex        int32  `json:"sex"`         // 性别,0:位置；1:男；2:女
	AvatarUrl  string `json:"img"`         // 用户头像
	UserExtra  string `json:"user_extra"`  // 用户附件属性
	MemberType int    `json:"member_type"` // 群组类型
	Remarks    string `json:"remarks"`     // 备注
	Extra      string `json:"extra"`       // 附件属性
}

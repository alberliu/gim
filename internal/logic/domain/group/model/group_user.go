package model

import "time"

const (
	UpdateTypeUpdate = 1
	UpdateTypeDelete = 2
)

type GroupUser struct {
	Id         int64     // 自增主键
	GroupId    int64     // 群组id
	UserId     int64     // 用户id
	MemberType int       // 群组类型
	Remarks    string    // 备注
	Extra      string    // 附加属性
	Status     int       // 状态
	CreateTime time.Time // 创建时间
	UpdateTime time.Time // 更新时间
	UpdateType int       `gorm:"-"` // 更新类型
}

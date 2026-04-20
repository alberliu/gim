package domain

import "time"

const (
	FriendStatusApply = 0 // 申请
	FriendStatusAgree = 1 // 同意
)

type Friend struct {
	UserID    uint64    // 用户ID
	FriendID  uint64    // 好友ID
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
	Remarks   string    // 备注
	Extra     string    // 扩展字段
	Status    int       // 状态
}

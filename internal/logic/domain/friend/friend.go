package friend

import "time"

const (
	FriendStatusApply = 0 // 申请
	FriendStatusAgree = 1 // 同意
)

type Friend struct {
	Id         int64
	UserId     int64
	FriendId   int64
	Remarks    string
	Extra      string
	Status     int
	CreateTime time.Time
	UpdateTime time.Time
}

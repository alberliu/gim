package domain

import (
	"time"

	pb "gim/pkg/protocol/pb/userpb"
)

// User 账户
type User struct {
	ID          uint64    // 用户id
	CreatedAt   time.Time // 创建时间
	UpdatedAt   time.Time // 更新时间
	PhoneNumber string    // 手机号
	Nickname    string    // 昵称
	Sex         int32     // 性别，1:男；2:女
	AvatarUrl   string    // 用户头像
	Extra       string    // 附加属性
}

func (u *User) ToProto() *pb.User {
	if u == nil {
		return nil
	}

	return &pb.User{
		UserId:     u.ID,
		Nickname:   u.Nickname,
		Sex:        u.Sex,
		AvatarUrl:  u.AvatarUrl,
		Extra:      u.Extra,
		CreateTime: u.CreatedAt.Unix(),
		UpdateTime: u.UpdatedAt.Unix(),
	}
}

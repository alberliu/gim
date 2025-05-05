package domain

import (
	"time"

	pb "gim/pkg/protocol/pb/logicpb"
)

type GroupUser struct {
	GroupID    uint64        // 群组id
	UserID     uint64        // 用户id
	CreatedAt  time.Time     // 创建时间
	UpdatedAt  time.Time     // 更新时间
	MemberType pb.MemberType // 成员类型
	Remarks    string        // 备注
	Extra      string        // 附加属性
	Status     int           // 状态
}

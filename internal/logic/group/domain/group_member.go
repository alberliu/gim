package domain

import (
	"time"

	pb "gim/pkg/protocol/pb/logicpb"
)

type GroupMember struct {
	GroupID   uint64               // 群组ID
	UserID    uint64               // 用户ID
	CreatedAt time.Time            // 创建时间
	UpdatedAt time.Time            // 更新时间
	Nickname  string               // 昵称
	Type      pb.GroupMemberType   // 类型
	Status    pb.GroupMemberStatus // 状态
	Extra     string               // 附加属性
}

func (m *GroupMember) TableName() string {
	return "group_member"
}

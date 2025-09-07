package domain

import (
	"time"

	pb "gim/pkg/protocol/pb/logicpb"
)

// Group 群组
type Group struct {
	ID           uint64    // 群组id
	CreatedAt    time.Time // 创建时间
	UpdatedAt    time.Time // 更新时间
	Name         string    // 组名
	AvatarUrl    string    // 头像
	Introduction string    // 群简介
	Extra        string    // 附加字段
	Members      []uint64  `gorm:"serializer:json"` // 群组成员
}

func (g *Group) ToProto() *pb.Group {
	if g == nil {
		return nil
	}

	return &pb.Group{
		Id:           g.ID,
		Name:         g.Name,
		AvatarUrl:    g.AvatarUrl,
		Introduction: g.Introduction,
		Extra:        g.Extra,
		Members:      g.Members,
	}
}

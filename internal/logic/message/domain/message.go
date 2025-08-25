package domain

import (
	"time"

	"gim/pkg/protocol/pb/connectpb"
)

type Message struct {
	ID        uint64            // 自增主键
	CreatedAt time.Time         // 创建时间
	UpdatedAt time.Time         // 更新时间
	RequestID string            // 请求id
	Command   connectpb.Command // 指令
	Content   []byte            // 消息内容
	Status    int8              // 消息状态，0：未处理;1：消息撤回
}

func (m *Message) TableName() string {
	return "message"
}

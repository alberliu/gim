package domain

import (
	"time"

	pb "gim/pkg/protocol/pb/logicpb"
)

type Message struct {
	ID        uint64      // 自增主键
	CreatedAt time.Time   // 创建时间
	UpdatedAt time.Time   // 更新时间
	RequestID int64       // 请求id
	Code      pb.PushCode // 推送码
	Content   []byte      // 消息内容
	Status    int8        // 消息状态，0：未处理;1：消息撤回
}

func (m *Message) TableName() string {
	return "message"
}

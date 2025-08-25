package domain

import (
	"time"

	"gim/pkg/protocol/pb/connectpb"
)

type UserMessage struct {
	UserID    uint64    // 用户ID
	Seq       uint64    // 序列号
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
	MessageID uint64    // 消息ID
	Message   Message   `gorm:"->"` // 消息
}

func (m *UserMessage) MessageToPB() *connectpb.Message {
	return &connectpb.Message{
		Command:   m.Message.Command,
		Content:   m.Message.Content,
		Seq:       m.Seq,
		CreatedAt: time.Now().Unix(),
	}
}

func MessagesToPB(messages []UserMessage) []*connectpb.Message {
	pbMessages := make([]*connectpb.Message, 0, len(messages))
	for i := range messages {
		pbMessage := messages[i].MessageToPB()
		if pbMessages != nil {
			pbMessages = append(pbMessages, pbMessage)
		}
	}
	return pbMessages
}

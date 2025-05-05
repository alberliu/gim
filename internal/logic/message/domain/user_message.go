package domain

import (
	"time"

	pb "gim/pkg/protocol/pb/logicpb"
	"gim/pkg/util"
)

type UserMessage struct {
	UserID    uint64    // 用户ID
	Seq       uint64    // 序列号
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
	MessageID uint64    // 消息ID
	Message   Message   `gorm:"->"` // 消息
}

func (m *UserMessage) MessageToPB() *pb.Message {
	return &pb.Message{
		Code:      m.Message.Code,
		Content:   m.Message.Content,
		Seq:       m.Seq,
		CreatedAt: util.UnixMilliTime(m.Message.CreatedAt),
		Status:    pb.MessageStatus(m.Message.Status),
	}
}

func MessagesToPB(messages []UserMessage) []*pb.Message {
	pbMessages := make([]*pb.Message, 0, len(messages))
	for i := range messages {
		pbMessage := messages[i].MessageToPB()
		if pbMessages != nil {
			pbMessages = append(pbMessages, pbMessage)
		}
	}
	return pbMessages
}

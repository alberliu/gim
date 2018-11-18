package model

import (
	"time"
)

// Message 消息
type Message struct {
	Id             int64     `json:"id"`               // 自增主键
	MessageId      int64     `json:"message_id"`       // 消息id
	UserId         int64     `json:"user_id"`          // 用户id
	SenderType     int       `json:"sender_type"`      // 发送者类型
	SenderId       int64     `json:"sender"`           // 发送者账户id
	SenderDeviceId int64     `json:"sender_device_id"` // 发送者设备id
	ReceiverType   int       `json:"receiver_type"`    // 接收者账户id
	ReceiverId     int64     `json:"receiver"`         // 接收者id,如果是单聊信息，则为user_id，如果是群组消息，则为group_id
	Type           int       `json:"type"`             // 消息类型,0：文本；1：语音；2：图片
	Content        string    `json:"content"`          // 内容
	Sequence       int64     `json:"sequence"`         // 消息同步序列
	SendTime       time.Time `json:"send_time"`        // 消息发送时间
	CreateTime     time.Time `json:"create_time"`      // 创建时间
}

package transfer

import (
	"goim/public/lib"
	"time"
)

const (
	MessageTypeSync = 1 // 消息同步
	MessageTypeMail = 2 // 消息投递
)

type Message struct {
	DeviceId int64         `json:"device_id"` // 设备id
	Type     int32         `json:"type"`      // 消息投递类型，1：消息同步列表，2：消息发送列表
	Messages []MessageItem `json:"messages"`  // 消息列表
}

// GetLog 获取消息日志
func (m *Message) GetLog() string {
	list := make([]messageLogItem, 0, len(m.Messages))
	var item messageLogItem
	for _, v := range m.Messages {
		item.MessageId = v.MessageId
		item.Sequence = v.Sequence
		list = append(list, item)
	}
	return lib.JsonMarshal(list)
}

type messageLogItem struct {
	MessageId int64 `json:"message_id"` // 消息id
	Sequence  int64 `json:"sequence"`   // 消息序列
}

// 单条消息投递
type MessageItem struct {
	MessageId      int64     `json:"message_id"`       // 消息id
	SenderType     int       `json:"sender_type"`      // 发送者类型
	SenderId       int64     `json:"sender_id"`        // 发送者id
	SenderDeviceId int64     `json:"sender_device_id"` // 发送者设备id
	ReceiverType   int       `json:"receiver_type"`    // 接收者类型
	ReceiverId     int64     `json:"receiver_id"`      // 接收者id
	Type           int       `json:"type"`             // 消息类型
	Content        string    `json:"content"`          // 消息内容
	Sequence       int64     `json:"sequence"`         // 消息序列
	SendTime       time.Time `json:"send_time"`        // 消息发送时间戳，精确到毫秒
}

type MessageACK struct {
	MessageId    int64     `json:"message_id"`    // 消息id
	DeviceId     int64     `json:"device_id"`     // 设备id
	UserId       int64     `json:"user_id"`       // 用户id
	SyncSequence int64     `json:"sync_sequence"` // 消息序列
	ReceiveTime  time.Time `json:"receive_time"`  // 消息接收时间戳，精确到毫秒
}

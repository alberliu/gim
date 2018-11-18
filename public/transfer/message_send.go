package transfer

import "time"

type MessageSend struct {
	MessageId      int64     `json:"message_id"`
	SenderDeviceId int64     `json:"sender_device_id"` // 发送者设备id
	SenderUserId   int64     `json:"sender_user_id"`   // 发送者用户id
	ReceiverType   int32     `json:"receiver_type"`    // 接收者类型，1：单发；2：群发
	ReceiverId     int64     `json:"receiver_id"`      // 接收者id
	Type           int32     `json:"type"`             // 消息类型
	Content        string    `json:"content"`          // 消息内容
	SendSequence   int64     `json:"send_sequence"`    // 消息序列号
	SendTime       time.Time `json:"send_time"`        // 消息发送时间戳，精确到毫秒
}

type MessageSendACK struct {
	MessageId    int64 `json:"message_id"`    // 消息id
	DeviceId     int64 `json:"device_id"`     // 设备id
	SendSequence int64 `json:"send_sequence"` // 消息序列号
	Code         int   `json:"code"`          // 消息发送结果
}

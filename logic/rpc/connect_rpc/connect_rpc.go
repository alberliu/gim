package connect_rpc

import (
	"goim/public/transfer"
)

// ConnectRPCer 连接层接口
type ConnectRPCer interface {
	SendMessage(message transfer.Message) error
	SendMessageSendACK(ack transfer.MessageSendACK) error
}

// ConnectRPCer 连接层实例
var ConnectRPC ConnectRPCer

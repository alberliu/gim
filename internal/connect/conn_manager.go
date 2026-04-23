package connect

import (
	"sync"

	pb "gim/pkg/protocol/pb/connectpb"
)

var ConnManager = sync.Map{}

// SetConn 存储
func SetConn(deviceID uint64, conn *Conn) {
	ConnManager.Store(deviceID, conn)
}

// GetConn 获取
func GetConn(deviceID uint64) *Conn {
	value, ok := ConnManager.Load(deviceID)
	if ok {
		return value.(*Conn)
	}
	return nil
}

// DeleteConn 删除
func DeleteConn(deviceID uint64) {
	ConnManager.Delete(deviceID)
}

// PushAll 全服推送
func PushAll(message *pb.Message) {
	ConnManager.Range(func(key, value any) bool {
		conn := value.(*Conn)
		conn.SendMessage(message)
		return true
	})
}

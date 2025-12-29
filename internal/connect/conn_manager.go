package connect

import (
	"sync"

	pb "gim/pkg/protocol/pb/connectpb"
)

var ConnesManager = sync.Map{}

// SetConn 存储
func SetConn(deviceID uint64, conn *Conn) {
	ConnesManager.Store(deviceID, conn)
}

// GetConn 获取
func GetConn(deviceID uint64) *Conn {
	value, ok := ConnesManager.Load(deviceID)
	if ok {
		return value.(*Conn)
	}
	return nil
}

// DeleteConn 删除
func DeleteConn(deviceID uint64) {
	ConnesManager.Delete(deviceID)
}

// PushAll 全服推送
func PushAll(message *pb.Message) {
	ConnesManager.Range(func(key, value interface{}) bool {
		conn := value.(*Conn)
		conn.SendMessage(message)
		return true
	})
}

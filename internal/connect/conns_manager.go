package connect

import (
	"sync"

	pb "gim/pkg/gen/proto/connectpb"
)

var ConnsManager = sync.Map{}

// SetConn 存储
func SetConn(deviceID uint64, conn *Conn) {
	ConnsManager.Store(deviceID, conn)
}

// GetConn 获取
func GetConn(deviceID uint64) *Conn {
	value, ok := ConnsManager.Load(deviceID)
	if ok {
		return value.(*Conn)
	}
	return nil
}

// DeleteConn 删除：仅当当前 map 中存储的就是 conn 本身时才删除，
// 防止旧连接的清理流程把已被新连接覆盖的条目误删。
func DeleteConn(deviceID uint64, conn *Conn) {
	ConnsManager.CompareAndDelete(deviceID, conn)
}

// PushAll 全服推送
func PushAll(message *pb.Message) {
	ConnsManager.Range(func(key, value any) bool {
		conn := value.(*Conn)
		conn.SendMessage(message)
		return true
	})
}

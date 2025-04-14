package connect

import (
	"sync"

	"gim/pkg/protocol/pb"
)

var ConnesManager = sync.Map{}

// SetConn 存储
func SetConn(deviceId int64, conn *Conn) {
	ConnesManager.Store(deviceId, conn)
}

// GetConn 获取
func GetConn(deviceId int64) *Conn {
	value, ok := ConnesManager.Load(deviceId)
	if ok {
		return value.(*Conn)
	}
	return nil
}

// DeleteConn 删除
func DeleteConn(deviceId int64) {
	ConnesManager.Delete(deviceId)
}

// PushAll 全服推送
func PushAll(message *pb.Message) {
	ConnesManager.Range(func(key, value interface{}) bool {
		conn := value.(*Conn)
		conn.Send(pb.PackageType_PT_MESSAGE, 0, message, nil)
		return true
	})
}

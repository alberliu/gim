package connect

import (
	"gim/pkg/pb"
	"sync"
)

var ConnsManager = sync.Map{}

// SetConn 存储
func SetConn(deviceId int64, conn *Conn) {
	ConnsManager.Store(deviceId, conn)
}

// GetConn 获取
func GetConn(deviceId int64) *Conn {
	value, ok := ConnsManager.Load(deviceId)
	if ok {
		return value.(*Conn)
	}
	return nil
}

// DeleteConn 删除
func DeleteConn(deviceId int64) {
	ConnsManager.Delete(deviceId)
}

func PushAll(message *pb.MessageSend) {
	ConnsManager.Range(func(key, value interface{}) bool {
		conn := value.(*Conn)
		conn.Send(pb.PackageType_PT_MESSAGE, 0, message, nil)
		return true
	})
}

package ws_conn

import "sync"

var manager sync.Map

// store 存储
func store(deviceId int64, ctx *WSConnContext) {
	manager.Store(deviceId, ctx)
}

// load 获取
func load(deviceId int64) *WSConnContext {
	value, ok := manager.Load(deviceId)
	if ok {
		return value.(*WSConnContext)
	}
	return nil
}

// delete 删除
func delete(deviceId int64) {
	manager.Delete(deviceId)
}

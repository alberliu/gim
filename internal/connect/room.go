package connect

import (
	"container/list"
	"log/slog"
	"sync"

	pb "gim/pkg/protocol/pb/connectpb"
)

var RoomsManager sync.Map

// SubscribedRoom 订阅房间
func SubscribedRoom(conn *Conn, roomID uint64) {
	if roomID == conn.RoomID {
		return
	}

	oldRoomID := conn.RoomID
	// 取消订阅
	if oldRoomID != 0 {
		value, ok := RoomsManager.Load(oldRoomID)
		if !ok {
			return
		}
		room := value.(*Room)
		room.Unsubscribe(conn)

		if room.Conns.Front() == nil {
			RoomsManager.Delete(oldRoomID)
		}
		slog.Debug("SubscribedRoom un", "userID", conn.UserID, "roomID", roomID)
		return
	}

	// 订阅
	if roomID != 0 {
		value, ok := RoomsManager.Load(roomID)
		var room *Room
		if !ok {
			room = NewRoom(roomID)
			RoomsManager.Store(roomID, room)
		} else {
			room = value.(*Room)
		}
		room.Subscribe(conn)
		slog.Debug("SubscribedRoom", "userID", conn.UserID, "roomID", roomID)
		return
	}
}

// PushRoom 房间消息推送
func PushRoom(roomID uint64, message *pb.Message) {
	value, ok := RoomsManager.Load(roomID)
	if !ok {
		return
	}

	slog.Debug("PushRoom", "roomID", roomID, "msg", message)
	value.(*Room).Push(message)
}

type Room struct {
	RoomID uint64     // 房间ID
	Conns  *list.List // 订阅房间消息的连接
	lock   sync.RWMutex
}

func NewRoom(roomID uint64) *Room {
	return &Room{
		RoomID: roomID,
		Conns:  list.New(),
	}
}

// Subscribe 订阅房间
func (r *Room) Subscribe(conn *Conn) {
	r.lock.Lock()
	defer r.lock.Unlock()

	conn.Element = r.Conns.PushBack(conn)
	conn.RoomID = r.RoomID
}

// Unsubscribe 取消订阅
func (r *Room) Unsubscribe(conn *Conn) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.Conns.Remove(conn.Element)
	conn.Element = nil
	conn.RoomID = 0
}

// Push 推送消息到房间
func (r *Room) Push(message *pb.Message) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	element := r.Conns.Front()
	for {
		conn := element.Value.(*Conn)
		slog.Debug("PushRoom toUser", "userID", conn.UserID, "msg", message)
		conn.SendMessage(message)

		element = element.Next()
		if element == nil {
			break
		}
	}
}

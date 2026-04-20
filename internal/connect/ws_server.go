package connect

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 65536,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// StartWSServer 启动WebSocket服务器
func StartWSServer(address string) {
	http.HandleFunc("/ws", wsHandler)
	slog.Info("websocket server running")
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("upgrade error", "error", err)
		return
	}
	DoConn(wsConn)
}

// DoConn 处理连接
func DoConn(wsConn *websocket.Conn) {
	conn := &Conn{
		ConnType: ConnTypeWS,
		WS:       wsConn,
	}

	for {
		err := conn.WS.SetReadDeadline(time.Now().Add(ReadDeadline))
		if err != nil {
			conn.Close(err)
			return
		}
		_, data, err := conn.WS.ReadMessage()
		if err != nil {
			conn.Close(err)
			return
		}

		conn.HandlePacket(data)
	}
}

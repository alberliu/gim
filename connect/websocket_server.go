package connect

import (
	"goim/public/logger"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type WebsocketServer struct {
}

func (*WebsocketServer) Start() {
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe("0.0.0.0:8888", nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Sugar.Error(err)
	}
	conn.WriteMessage(1, []byte("hello"))
}

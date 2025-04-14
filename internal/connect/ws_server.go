package connect

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"gim/pkg/logger"
	"gim/pkg/util"
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
	logger.Logger.Info("websocket server start")
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	conn := &Conn{
		CoonType: ConnTypeWS,
		WS:       wsConn,
	}
	DoConn(conn)
}

// DoConn 处理连接
func DoConn(conn *Conn) {
	defer util.RecoverPanic()

	for {
		err := conn.WS.SetReadDeadline(time.Now().Add(12 * time.Minute))
		if err != nil {
			HandleReadErr(conn, err)
			return
		}
		_, data, err := conn.WS.ReadMessage()
		if err != nil {
			HandleReadErr(conn, err)
			return
		}

		conn.HandleMessage(data)
	}
}

// HandleReadErr 读取conn错误
func HandleReadErr(conn *Conn, err error) {
	logger.Logger.Debug("read tcp error：", zap.Int64("user_id", conn.UserId),
		zap.Int64("device_id", conn.DeviceId), zap.Error(err))
	str := err.Error()
	// 服务器主动关闭连接
	if strings.HasSuffix(str, "use of closed network connection") {
		return
	}

	conn.Close()
	// 客户端主动关闭连接或者异常程序退出
	if err == io.EOF {
		return
	}
	// SetReadDeadline 之后，超时返回的错误
	if strings.HasSuffix(str, "i/o timeout") {
		return
	}
}

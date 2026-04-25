package connect

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"gim/pkg/md"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/rpc"
	"gim/pkg/safe"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 65536,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// StartWSServer 启动WebSocket服务器，返回 server 供调用方平滑关闭
func StartWSServer(address string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHandler)
	server := &http.Server{Addr: address, Handler: mux}
	go func() {
		slog.Info("websocket server running")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	return server
}

func parseRequest(r *http.Request) (*logicpb.SignInRequest, error) {
	q := r.URL.Query()

	userID, err := strconv.ParseUint(q.Get("user_id"), 10, 64)
	if err != nil {
		return nil, err
	}
	deviceID, err := strconv.ParseUint(q.Get("device_id"), 10, 64)
	if err != nil {
		return nil, err
	}
	token := q.Get("token")
	if token == "" {
		return nil, errors.New("token is empty")
	}

	clientAddr := r.Header.Get("X-Real-IP")
	if clientAddr == "" {
		clientAddr = r.RemoteAddr
	}
	return &logicpb.SignInRequest{
		UserId:     userID,
		DeviceId:   deviceID,
		Token:      token,
		ClientAddr: clientAddr,
	}, nil
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	defer safe.RecoverPanic()

	request, err := parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = rpc.GetDeviceIntClient().SignIn(md.WithGenerateRequestID(), request)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("upgrade websocket failed", "error", err)
		return
	}

	conn := &Conn{
		ConnType: ConnTypeWS,
		WS:       wsConn,
		UserID:   request.UserId,
		DeviceID: request.DeviceId,
	}
	SetConn(request.DeviceId, conn)
	handleMessage(conn)
}

// handleMessage 处理消息
func handleMessage(conn *Conn) {
	var err error
	var buf []byte
	defer func() { conn.Close(err) }()

	for {
		err = conn.WS.SetReadDeadline(time.Now().Add(ReadDeadline))
		if err != nil {
			return
		}
		_, buf, err = conn.WS.ReadMessage()
		if err != nil {
			return
		}

		conn.HandlePacket(buf)
	}
}

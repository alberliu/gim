package test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"gim/pkg/codec"
	"gim/pkg/protocol/pb/businesspb"
	"gim/pkg/protocol/pb/connectpb"
	"gim/pkg/protocol/pb/logicpb"
)

var (
	connectTcpServerAddr = "127.0.0.1:8001"
	businessServerAddr   = "127.0.0.1:8020"
	logicServerAddr      = "127.0.0.1:8010"
)

func connect(userID, deviceID uint64) {
	log := slog.With("userID", userID, "deviceID", deviceID)

	reply, err := getUserExtServiceClient().SignIn(context.TODO(), &businesspb.SignInRequest{
		PhoneNumber: strconv.FormatUint(userID, 10),
		Code:        "0",
		Device: &logicpb.Device{
			Id:            deviceID,
			Type:          logicpb.DeviceType_DT_ANDROID,
			Brand:         "xiaomi",
			Model:         "xiaomi 15",
			SystemVersion: "15.0.0",
			SdkVersion:    "1.0.0",
			BranchPushId:  "xiaomi push id",
		},
	})
	if err != nil {
		panic(err)
	}
	log.Info("短连接登录成功", "reply", reply)
	go runClient(log, "tcp", connectTcpServerAddr, userID, deviceID, reply.Token)
}

func getUserExtServiceClient() businesspb.UserExtServiceClient {
	conn, err := grpc.NewClient(businessServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return businesspb.NewUserExtServiceClient(conn)
}

type conn interface {
	write(buf []byte) error
	receive(handler func([]byte))
}

type tcpConn struct {
	conn   net.Conn
	reader *bufio.Reader
}

func newTCPConn(url string) (*tcpConn, error) {
	// demo "127.0.0.1:8001"
	conn, err := net.Dial("tcp", url)
	if err != nil {
		return nil, err
	}

	return &tcpConn{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}, nil
}

func (c *tcpConn) write(buf []byte) error {
	_, err := c.conn.Write(codec.Encode(buf))
	return err
}

func (c *tcpConn) receive(handler func([]byte)) {
	for {
		buf, err := codec.Decode(c.reader)
		if err != nil {
			log.Println(err)
			return
		}

		handler(buf)
	}
}

type wsConn struct {
	conn *websocket.Conn
}

func newWsConn(url string) (*wsConn, error) {
	// demo "ws://127.0.0.1:8002/ws"
	conn, resp, err := websocket.DefaultDialer.Dial(url, http.Header{})
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(string(bytes))
	return &wsConn{conn: conn}, nil
}

func (c *wsConn) write(buf []byte) error {
	return c.conn.WriteMessage(websocket.BinaryMessage, buf)
}

func (c *wsConn) receive(handler func([]byte)) {
	for {
		_, bytes, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		handler(bytes)
	}
}

type client struct {
	UserID   uint64
	DeviceID uint64
	Token    string
	conn     conn
	log      *slog.Logger
}

func runClient(log *slog.Logger, network string, addr string, userID, deviceID uint64, token string) {
	var conn conn
	var err error
	switch network {
	case "tcp":
		conn, err = newTCPConn(addr)
	case "ws":
		conn, err = newWsConn(addr)
	default:
		panic("unsupported network")
	}
	if err != nil {
		panic(err)
	}

	client := &client{
		UserID:   userID,
		DeviceID: deviceID,
		Token:    token,
		conn:     conn,
		log:      log,
	}
	client.run()
}

func (c *client) run() {
	go c.conn.receive(c.handleMessage)
	c.signIn()
	c.subscribeRoom()
	c.heartbeat()
}

func (c *client) send(pt connectpb.Command, requestID string, msg proto.Message) {
	var message = connectpb.Message{
		Command:   pt,
		RequestId: requestID,
	}

	if msg != nil {
		bytes, err := proto.Marshal(msg)
		if err != nil {
			c.log.Error("send", "error", err)
			return
		}
		message.Content = bytes
	}

	buf, err := proto.Marshal(&message)
	if err != nil {
		c.log.Error("send", "error", err)
		return
	}

	err = c.conn.write(buf)
	if err != nil {
		c.log.Error("send", "error", err)
	}
}

func getRequestID() string {
	unix := time.Now().UnixNano()
	return strconv.FormatInt(unix, 10)
}

func (c *client) signIn() {
	request := logicpb.SignInRequest{
		UserId:   c.UserID,
		DeviceId: c.DeviceID,
		Token:    c.Token,
	}
	c.send(connectpb.Command_SIGN_IN, getRequestID(), &request)
	c.log.Info("发送登录指令")
	time.Sleep(1 * time.Second)
}

func (c *client) heartbeat() {
	ticker := time.NewTicker(time.Minute * 5)
	for range ticker.C {
		c.send(connectpb.Command_HEARTBEAT, getRequestID(), nil)
		c.log.Info("心跳发送")
	}
}

func (c *client) subscribeRoom() {
	var roomID uint64 = 1
	c.send(connectpb.Command_SUBSCRIBE_ROOM, getRequestID(), &logicpb.SubscribeRoomRequest{
		RoomId: roomID,
	})
	c.log.Info("订阅房间", "roomID", roomID)
}

func (c *client) handleMessage(buf []byte) {
	var message connectpb.Message
	err := proto.Unmarshal(buf, &message)
	if err != nil {
		log.Println(err)
		return
	}

	switch message.Command {
	case connectpb.Command_SIGN_IN:
		c.log.Info("登录响应", "message", jsonString(&message), "reply", jsonString(getReply(&message)))

		time.Sleep(1 * time.Second)
	case connectpb.Command_HEARTBEAT:
		c.log.Info("心跳响应")
	case connectpb.Command_SUBSCRIBE_ROOM:
		c.log.Info("订阅房间响应", "message", jsonString(&message), "reply", jsonString(getReply(&message)))
	default:
		c.log.Info("other", "message", &message)
	}
}

func getReply(message *connectpb.Message) *connectpb.Reply {
	var reply connectpb.Reply
	_ = proto.Unmarshal(message.Content, &reply)
	return &reply
}

func jsonString(any any) string {
	bytes, _ := json.Marshal(any)
	return string(bytes)
}

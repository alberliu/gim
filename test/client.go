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
	connectTCPServerAddr       = "127.0.0.1:8002"
	connectWebsocketServerAddr = "ws://127.0.0.1:8003/ws"
	businessServerAddr         = "127.0.0.1:8020"
	logicServerAddr            = "127.0.0.1:8010"
)

type Network string

const (
	NetworkTCP       Network = "tcp"
	NetworkWebsocket Network = "websocket"
)

func getUserExtServiceClient() businesspb.UserExtServiceClient {
	conn, err := grpc.NewClient(businessServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return businesspb.NewUserExtServiceClient(conn)
}

func getMessageIntClient() logicpb.MessageIntServiceClient {
	conn, err := grpc.NewClient(logicServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return logicpb.NewMessageIntServiceClient(conn)
}

func getGroupIntClient() logicpb.GroupIntServiceClient {
	conn, err := grpc.NewClient(logicServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return logicpb.NewGroupIntServiceClient(conn)
}

func getRoomIntClient() logicpb.RoomIntServiceClient {
	conn, err := grpc.NewClient(logicServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return logicpb.NewRoomIntServiceClient(conn)
}

func connect(network Network, userID, deviceID uint64) {
	log := slog.With("userID", userID, "deviceID", deviceID)

	reply, err := getUserExtServiceClient().SignIn(context.Background(), &businesspb.SignInRequest{
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
	log.Info("sign in success", "reply", reply)

	addr := connectTCPServerAddr
	if network == NetworkWebsocket {
		addr = connectWebsocketServerAddr
	}
	go runClient(log, network, addr, userID, deviceID, reply.Token)
}

type client struct {
	Network  Network
	UserID   uint64
	DeviceID uint64
	Token    string
	conn     conn
	log      *slog.Logger
}

func runClient(log *slog.Logger, network Network, addr string, userID, deviceID uint64, token string) {
	var conn conn
	var err error
	switch network {
	case NetworkTCP:
		conn, err = newTCPConn(addr)
	case NetworkWebsocket:
		conn, err = newWsConn(addr, userID, deviceID, token)
	default:
		panic("unsupported network")
	}
	if err != nil {
		panic(err)
	}

	client := &client{
		Network:  network,
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

func (c *client) send(pt connectpb.PacketCommand, requestID string, msg proto.Message) {
	var packet = connectpb.Packet{
		Command:   pt,
		RequestId: requestID,
	}

	if msg != nil {
		bytes, err := proto.Marshal(msg)
		if err != nil {
			c.log.Error("marshal message failed", "error", err)
			return
		}
		packet.Content = bytes
	}

	buf, err := proto.Marshal(&packet)
	if err != nil {
		c.log.Error("marshal packet failed", "error", err)
		return
	}

	err = c.conn.write(buf)
	if err != nil {
		c.log.Error("write packet failed", "error", err)
	}
}

func (c *client) signIn() {
	if c.Network == NetworkWebsocket {
		return
	}

	request := logicpb.SignInRequest{
		UserId:   c.UserID,
		DeviceId: c.DeviceID,
		Token:    c.Token,
	}
	c.send(connectpb.PacketCommand_PC_SIGN_IN, generateRequestID(), &request)
	c.log.Info("send sign in")
	time.Sleep(1 * time.Second)
}

func (c *client) heartbeat() {
	ticker := time.NewTicker(time.Minute * 5)
	for range ticker.C {
		c.send(connectpb.PacketCommand_PC_HEARTBEAT, generateRequestID(), nil)
		c.log.Info("send heartbeat")
	}
}

func (c *client) subscribeRoom() {
	var roomID uint64 = 1
	c.send(connectpb.PacketCommand_PC_SUBSCRIBE_ROOM, generateRequestID(), &connectpb.SubscribeRoomRequest{
		RoomId: roomID,
	})
	c.log.Info("subscribe room", "roomID", roomID)
}

func (c *client) handleMessage(buf []byte) {
	var packet connectpb.Packet
	err := proto.Unmarshal(buf, &packet)
	if err != nil {
		log.Println(err)
		return
	}

	switch packet.Command {
	case connectpb.PacketCommand_PC_SIGN_IN:
		c.log.Info("sign in response", "packet", jsonString(&packet), "reply", jsonString(getReply(&packet)))

		time.Sleep(1 * time.Second)
	case connectpb.PacketCommand_PC_HEARTBEAT:
		c.log.Info("heartbeat response")
	case connectpb.PacketCommand_PC_SUBSCRIBE_ROOM:
		c.log.Info("subscribe room response", "packet", jsonString(&packet), "reply", jsonString(getReply(&packet)))
	case connectpb.PacketCommand_PC_MESSAGE:
		var message connectpb.Message
		err := proto.Unmarshal(packet.Content, &message)
		if err != nil {
			c.log.Error("unmarshal message failed", "error", err)
			return
		}
		c.log.Info("receive message", "message", stringMessage(&message))
	default:
		c.log.Info("other", "packet", &packet)
	}
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

func newWsConn(url string, userID, deviceID uint64, token string) (*wsConn, error) {
	// demo "ws://127.0.0.1:8002/ws"
	url = fmt.Sprintf("%s?user_id=%d&device_id=%d&token=%s", url, userID, deviceID, token)
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

func generateRequestID() string {
	unix := time.Now().UnixNano()
	return strconv.FormatInt(unix, 10)
}

func stringMessage(message *connectpb.Message) string {
	return jsonString(map[string]any{
		"command": message.Command,
		"content": string(message.Content),
	})
}

func getReply(packet *connectpb.Packet) *connectpb.Reply {
	var reply connectpb.Reply
	_ = proto.Unmarshal(packet.Content, &reply)
	return &reply
}

func jsonString(any any) string {
	bytes, _ := json.Marshal(any)
	return string(bytes)
}

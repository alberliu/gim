package connect

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	"gim/pkg/codec"
	pb "gim/pkg/protocol/pb/connectpb"
)

func TestTCPClient(t *testing.T) {
	runClient("tcp", "127.0.0.1:8001", 1, 1)
}

func TestWSClient(t *testing.T) {
	runClient("ws", "ws://127.0.0.1:8002/ws", 1, 1)
}

func TestGroupTCPClient(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	go runClient("tcp", "127.0.0.1:8001", 1, 1)
	go runClient("tcp", "127.0.0.1:8001", 2, 2)
	select {}
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
	conn     conn
}

func runClient(network string, url string, userID, deviceID uint64) {
	var conn conn
	var err error
	switch network {
	case "tcp":
		conn, err = newTCPConn(url)
	case "ws":
		conn, err = newWsConn(url)
	default:
		panic("unsupported network")
	}
	if err != nil {
		panic(err)
	}

	client := &client{
		UserID:   userID,
		DeviceID: deviceID,
		conn:     conn,
	}
	client.run()
}

func (c *client) run() {
	go c.conn.receive(c.handleMessage)
	c.signIn()
	c.subscribeRoom()
	c.heartbeat()
}

func (c *client) info() string {
	return fmt.Sprintf("%-5d%-5d", c.UserID, c.DeviceID)
}

func (c *client) send(pt pb.Command, requestID string, msg proto.Message) {
	var message = pb.Message{
		Command:   pt,
		RequestId: requestID,
	}

	if msg != nil {
		bytes, err := proto.Marshal(msg)
		if err != nil {
			log.Println(c.info(), err)
			return
		}
		message.Content = bytes
	}

	buf, err := proto.Marshal(&message)
	if err != nil {
		log.Println(c.info(), err)
		return
	}

	err = c.conn.write(buf)
	if err != nil {
		log.Println(c.info(), err)
	}
}

func getRequestID() string {
	unix := time.Now().UnixNano()
	return strconv.FormatInt(unix, 10)
}

func (c *client) signIn() {
	request := pb.SignInRequest{
		UserId:   c.UserID,
		DeviceId: 1,
		Token:    "0",
	}
	c.send(pb.Command_SIGN_IN, getRequestID(), &request)
	log.Println(c.info(), "发送登录指令")
	time.Sleep(1 * time.Second)
}

func (c *client) heartbeat() {
	ticker := time.NewTicker(time.Minute * 5)
	for range ticker.C {
		c.send(pb.Command_HEARTBEAT, getRequestID(), nil)
		fmt.Println(c.info(), "心跳发送")
	}
}

func (c *client) subscribeRoom() {
	var roomID uint64 = 1
	c.send(pb.Command_SUBSCRIBE_ROOM, getRequestID(), &pb.SubscribeRoomRequest{
		RoomId: roomID,
	})
	log.Println(c.info(), "订阅房间:", roomID)
}

func (c *client) handleMessage(buf []byte) {
	var message pb.Message
	err := proto.Unmarshal(buf, &message)
	if err != nil {
		log.Println(err)
		return
	}

	switch message.Command {
	case pb.Command_SIGN_IN:
		log.Println(c.info(), "登录响应:", jsonString(&message), jsonString(getReply(&message)))

		time.Sleep(1 * time.Second)
	case pb.Command_HEARTBEAT:
		log.Println(c.info(), "心跳响应")
	case pb.Command_SUBSCRIBE_ROOM:
		log.Println(c.info(), "订阅房间响应", jsonString(&message), jsonString(getReply(&message)))
	default:
		log.Println(c.info(), "other", &message)
	}
}

func getReply(message *pb.Message) *pb.Reply {
	var reply pb.Reply
	_ = proto.Unmarshal(message.Content, &reply)
	return &reply
}

func jsonString(any any) string {
	bytes, _ := json.Marshal(any)
	return string(bytes)
}

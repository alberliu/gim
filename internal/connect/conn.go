package connect

import (
	"bufio"
	"container/list"
	"context"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"gim/config"
	"gim/pkg/codec"
	"gim/pkg/gerrors"
	"gim/pkg/md"
	pb "gim/pkg/protocol/pb/connectpb"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/rpc"
)

const WriteDeadline = time.Second * 3

const (
	ConnTypeTCP int8 = 1 // tcp连接
	ConnTypeWS  int8 = 2 // websocket连接
)

type Conn struct {
	ConnType int8 // 连接类型

	TCP    *net.TCPConn  // tcp连接
	Reader *bufio.Reader // reader

	WSMutex sync.Mutex      // WS写锁
	WS      *websocket.Conn // websocket连接

	UserID   uint64        // 用户ID
	DeviceID uint64        // 设备ID
	RoomID   uint64        // 订阅的房间ID
	Element  *list.Element // 链表节点
}

// Write 写入数据
func (c *Conn) Write(buf []byte) error {
	var err error
	switch c.ConnType {
	case ConnTypeTCP:
		err = c.WriteToTCP(buf)
	case ConnTypeWS:
		err = c.WriteToWS(buf)
	}

	if err != nil {
		c.Close(err)
	}
	return err
}

// WriteToTCP 消息写入TCP
func (c *Conn) WriteToTCP(buf []byte) error {
	err := c.TCP.SetWriteDeadline(time.Now().Add(WriteDeadline))
	if err != nil {
		return err
	}

	_, err = c.TCP.Write(codec.Encode(buf))
	return err
}

// WriteToWS 消息写入WebSocket
func (c *Conn) WriteToWS(buf []byte) error {
	c.WSMutex.Lock()
	defer c.WSMutex.Unlock()

	err := c.WS.SetWriteDeadline(time.Now().Add(WriteDeadline))
	if err != nil {
		return err
	}
	return c.WS.WriteMessage(websocket.BinaryMessage, buf)
}

// Close 关闭
// use of closed network connection 服务端主动关闭
// io.EOF是用户主动断开连接
// io timeout是SetReadDeadline之后，超时返回的错误
func (c *Conn) Close(err error) {
	// 取消设备和连接的对应关系
	if c.DeviceID != 0 {
		DeleteConn(c.DeviceID)
	}

	// 取消订阅，需要异步出去，防止重复加锁造成死锁
	go func() {
		SubscribedRoom(c, 0)
	}()

	if c.DeviceID != 0 {
		_, _ = rpc.GetDeviceIntClient().Offline(context.TODO(), &logicpb.OfflineRequest{
			UserId:     c.UserID,
			DeviceId:   c.DeviceID,
			ClientAddr: c.GetAddr(),
		})
	}

	switch c.ConnType {
	case ConnTypeTCP:
		_ = c.TCP.Close()
	case ConnTypeWS:
		_ = c.WS.Close()
	}
}

func (c *Conn) GetAddr() string {
	switch c.ConnType {
	case ConnTypeTCP:
		return c.TCP.RemoteAddr().String()
	case ConnTypeWS:
		return c.WS.RemoteAddr().String()
	}
	return ""
}

// HandleMessage 消息处理
func (c *Conn) HandleMessage(buf []byte) {
	var message = new(pb.Message)
	err := proto.Unmarshal(buf, message)
	if err != nil {
		slog.Error("unmarshal error", "error", err, "len", len(buf))
		return
	}
	slog.Debug("HandleMessage", "message", message)

	// 对未登录的用户进行拦截
	if message.Command != pb.Command_SIGN_IN && c.UserID == 0 {
		setContent(message, gerrors.ErrUnauthorized, nil)
		c.Send(message)
		return
	}

	switch message.Command {
	case pb.Command_SIGN_IN:
		c.SignIn(message)
	case pb.Command_HEARTBEAT:
		c.Heartbeat(message)
	case pb.Command_SUBSCRIBE_ROOM:
		c.SubscribedRoom(message)
	default:
		slog.Error("handler switch other", "command", message.Command)
	}
}

// Send 下发消息
func (c *Conn) Send(message *pb.Message) {
	buf, err := proto.Marshal(message)
	if err != nil {
		slog.Error("proto.Marshal error", "error", err)
		return
	}

	err = c.Write(buf)
	if err != nil {
		slog.Error("Write error", "error", err)
		c.Close(err)
		return
	}
	slog.Info("Send", "userID", c.UserID, "message", message)
}

// SignIn 登录
func (c *Conn) SignIn(message *pb.Message) {
	var request pb.SignInRequest
	err := proto.Unmarshal(message.Content, &request)
	if err != nil {
		slog.Error("proto unmarshal error", "error", err)
		return
	}

	_, err = rpc.GetDeviceIntClient().SignIn(md.ContextWithRequestID(context.TODO(), message.RequestId), &logicpb.SignInRequest{
		UserId:      request.UserId,
		DeviceId:    request.DeviceId,
		Token:       request.Token,
		ConnectAddr: config.Config.ConnectLocalAddr,
		ClientAddr:  c.GetAddr(),
	})

	setContent(message, err, nil)
	c.Send(message)
	if err != nil {
		return
	}

	c.UserID = request.UserId
	c.DeviceID = request.DeviceId
	SetConn(request.DeviceId, c)
}

// Heartbeat 心跳
func (c *Conn) Heartbeat(message *pb.Message) {
	c.Send(message)

	_, err := rpc.GetDeviceIntClient().Heartbeat(context.TODO(), &logicpb.HeartbeatRequest{
		UserId:   c.UserID,
		DeviceId: c.DeviceID,
	})
	if err != nil {
		slog.Error("Heartbeat error", "deviceID", c.DeviceID, "userID", c.UserID, "error", err)
	}

	slog.Info("heartbeat", "deviceID", c.DeviceID, "userID", c.UserID)
}

// SubscribedRoom 订阅房间
func (c *Conn) SubscribedRoom(message *pb.Message) {
	var subscribeRoom pb.SubscribeRoomRequest
	err := proto.Unmarshal(message.Content, &subscribeRoom)
	if err != nil {
		slog.Error("proto unmarshal", "error", err)
		return
	}

	SubscribedRoom(c, subscribeRoom.RoomId)
	setContent(message, nil, nil)
	c.Send(message)
	_, err = rpc.GetRoomIntClient().SubscribeRoom(context.TODO(), &logicpb.SubscribeRoomRequest{
		UserId:   c.UserID,
		DeviceId: c.DeviceID,
		RoomId:   subscribeRoom.RoomId,
		ConnAddr: config.Config.ConnectLocalAddr,
	})
	if err != nil {
		slog.Error("SubscribedRoom error", "error", err)
	}
}

func setContent(message *pb.Message, err error, data proto.Message) {
	var reply pb.Reply
	if err != nil {
		statusErr, _ := status.FromError(err)
		if statusErr == nil {
			return
		}

		reply.Code = int32(statusErr.Code())
		reply.Message = statusErr.Message()
	}

	if data != nil {
		reply.Data, err = proto.Marshal(data)
		if err != nil {
			slog.Error("setContent error", "error", err)
		}
		return
	}

	buf, err := proto.Marshal(&reply)
	if err != nil {
		slog.Error("setContent error", "error", err)
		return
	}
	message.Content = buf
}

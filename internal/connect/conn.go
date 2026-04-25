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

	"gim/pkg/codec"
	"gim/pkg/gerrors"
	"gim/pkg/md"
	pb "gim/pkg/protocol/pb/connectpb"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/rpc"
)

const ReadDeadline = time.Minute * 12
const WriteDeadline = time.Second * 5

type ConnType int

const (
	ConnTypeTCP ConnType = 1 // TCP连接
	ConnTypeWS  ConnType = 2 // WebSocket连接
)

type Conn struct {
	ConnType ConnType   // 连接类型
	Mutex    sync.Mutex // WS写锁

	TCP    *net.TCPConn  // TCP连接
	Reader *bufio.Reader // Reader

	WS *websocket.Conn // WebSocket连接

	UserID   uint64        // 用户ID
	DeviceID uint64        // 设备ID
	RoomID   uint64        // 订阅的房间ID
	Element  *list.Element // 链表节点
}

// Write 写入数据
func (c *Conn) Write(buf []byte) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	var err error
	switch c.ConnType {
	case ConnTypeTCP:
		err = c.WriteToTCP(buf)
	case ConnTypeWS:
		err = c.WriteToWS(buf)
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
	err := c.WS.SetWriteDeadline(time.Now().Add(WriteDeadline))
	if err != nil {
		return err
	}
	return c.WS.WriteMessage(websocket.BinaryMessage, buf)
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

// HandlePacket 包处理
func (c *Conn) HandlePacket(buf []byte) {
	packet := new(pb.Packet)
	err := proto.Unmarshal(buf, packet)
	if err != nil {
		slog.Error("unmarshal packet failed", "error", err, "len", len(buf))
		return
	}
	slog.Debug("handle packet", "packet", packet)

	// 对未登录的用户进行拦截
	if packet.Command != pb.PacketCommand_PC_SIGN_IN && c.UserID == 0 {
		setContent(packet, gerrors.ErrUnauthorized, nil)
		c.SendPacket(packet)
		return
	}

	switch packet.Command {
	case pb.PacketCommand_PC_SIGN_IN:
		c.SignIn(packet)
	case pb.PacketCommand_PC_HEARTBEAT:
		c.Heartbeat(packet)
	case pb.PacketCommand_PC_SUBSCRIBE_ROOM:
		c.SubscribedRoom(packet)
	default:
		slog.Error("unknown command", "command", packet.Command)
	}
}

// SendPacket 下发包
func (c *Conn) SendPacket(packet *pb.Packet) {
	buf, err := proto.Marshal(packet)
	if err != nil {
		slog.Error("marshal packet failed", "error", err)
		return
	}

	err = c.Write(buf)
	if err != nil {
		slog.Error("write packet failed", "error", err)
		c.Close(err)
		return
	}
	slog.Info("send packet", "user_id", c.UserID, "packet", packet)
}

func (c *Conn) SendMessage(message *pb.Message) {
	buf, err := proto.Marshal(message)
	if err != nil {
		slog.Error("marshal message failed", "error", err)
		return
	}

	packet := &pb.Packet{
		Command: pb.PacketCommand_PC_MESSAGE,
		Content: buf,
	}
	c.SendPacket(packet)

	slog.Info("send message", "user_id", c.UserID, "message", message)
}

// SignIn 登录
func (c *Conn) SignIn(packet *pb.Packet) {
	if c.ConnType != ConnTypeTCP {
		return
	}

	var request pb.SignInRequest
	err := proto.Unmarshal(packet.Content, &request)
	if err != nil {
		slog.Error("unmarshal sign in failed", "error", err)
		return
	}

	ctx, cancel := context.WithTimeout(md.WithRequestID(context.Background(), packet.RequestId), rpc.Timeout)
	defer cancel()
	_, err = rpc.GetDeviceIntClient().SignIn(ctx, &logicpb.SignInRequest{
		UserId:     request.UserId,
		DeviceId:   request.DeviceId,
		Token:      request.Token,
		ClientAddr: c.GetAddr(),
	})

	setContent(packet, err, nil)
	c.SendPacket(packet)
	if err != nil {
		return
	}

	c.UserID = request.UserId
	c.DeviceID = request.DeviceId
	SetConn(request.DeviceId, c)
}

// Heartbeat 心跳
func (c *Conn) Heartbeat(packet *pb.Packet) {
	c.SendPacket(packet)

	ctx, cancel := context.WithTimeout(context.Background(), rpc.Timeout)
	defer cancel()
	_, err := rpc.GetDeviceIntClient().Heartbeat(ctx, &logicpb.HeartbeatRequest{
		UserId:   c.UserID,
		DeviceId: c.DeviceID,
	})
	if err != nil {
		slog.Error("heartbeat failed", "device_id", c.DeviceID, "user_id", c.UserID, "error", err)
	}

	slog.Info("heartbeat received", "device_id", c.DeviceID, "user_id", c.UserID)
}

// SubscribedRoom 订阅房间
func (c *Conn) SubscribedRoom(packet *pb.Packet) {
	var subscribeRoom pb.SubscribeRoomRequest
	err := proto.Unmarshal(packet.Content, &subscribeRoom)
	if err != nil {
		slog.Error("unmarshal subscribe room failed", "error", err)
		return
	}

	SubscribedRoom(c, subscribeRoom.RoomId)
	setContent(packet, nil, nil)
	c.SendPacket(packet)
	ctx, cancel := context.WithTimeout(context.Background(), rpc.Timeout)
	defer cancel()
	_, err = rpc.GetRoomIntClient().SubscribeRoom(ctx, &logicpb.SubscribeRoomRequest{
		UserId:   c.UserID,
		DeviceId: c.DeviceID,
		RoomId:   subscribeRoom.RoomId,
	})
	if err != nil {
		slog.Error("subscribe room failed", "error", err)
	}
}

func setContent(packet *pb.Packet, err error, message proto.Message) {
	var reply pb.Reply
	if err != nil {
		statusErr := status.Convert(err)
		reply.Code = int32(statusErr.Code())
		reply.Message = statusErr.Message()
	}

	if message != nil {
		reply.Data, err = proto.Marshal(message)
		if err != nil {
			slog.Error("marshal reply data failed", "error", err)
		}
	}

	buf, err := proto.Marshal(&reply)
	if err != nil {
		slog.Error("marshal reply failed", "error", err)
		return
	}
	packet.Content = buf
}

// Close 关闭
// use of closed network connection 服务端主动关闭
// io.EOF是用户主动断开连接
// io timeout是SetReadDeadline之后，超时返回的错误
func (c *Conn) Close(err error) {
	slog.Warn("close conn", "error", err)
	// 取消设备和连接的对应关系
	if c.DeviceID != 0 {
		DeleteConn(c.DeviceID)
	}

	// 取消订阅，需要异步出去，防止重复加锁造成死锁
	go func() {
		SubscribedRoom(c, 0)
	}()

	if c.DeviceID != 0 {
		ctx, cancel := context.WithTimeout(context.Background(), rpc.Timeout)
		defer cancel()
		_, _ = rpc.GetDeviceIntClient().Offline(ctx, &logicpb.OfflineRequest{
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

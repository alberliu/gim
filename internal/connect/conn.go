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
	"gim/pkg/md"
	pb "gim/pkg/protocol/pb/connectpb"
	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/rpc"
)

const (
	CoonTypeTCP int8 = 1 // tcp连接
	ConnTypeWS  int8 = 2 // websocket连接
)

type Conn struct {
	CoonType int8 // 连接类型

	TCP    *net.TCPConn  // tcp连接
	Reader *bufio.Reader // reader

	WSMutex sync.Mutex      // WS写锁
	WS      *websocket.Conn // websocket连接

	UserId   uint64        // 用户ID
	DeviceId uint64        // 设备ID
	RoomId   uint64        // 订阅的房间ID
	Element  *list.Element // 链表节点
}

// Write 写入数据
func (c *Conn) Write(buf []byte) error {
	var err error
	switch c.CoonType {
	case CoonTypeTCP:
		err = c.WriteToTCP(buf)
	case ConnTypeWS:
		err = c.WriteToWS(buf)
	}

	if err != nil {
		c.Close(err)
	}
	return err
}

// WriteToTCP 消息写入WebSocket
func (c *Conn) WriteToTCP(buf []byte) error {
	err := c.TCP.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))
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

	err := c.WS.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))
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
	if c.DeviceId != 0 {
		DeleteConn(c.DeviceId)
	}

	// 取消订阅，需要异步出去，防止重复加锁造成死锁
	go func() {
		SubscribedRoom(c, 0)
	}()

	if c.DeviceId != 0 {
		_, _ = rpc.GetDeviceIntClient().Offline(context.TODO(), &logicpb.OfflineRequest{
			UserId:     c.UserId,
			DeviceId:   c.DeviceId,
			ClientAddr: c.GetAddr(),
		})
	}

	if c.CoonType == CoonTypeTCP {
		c.TCP.Close()
	} else if c.CoonType == ConnTypeWS {
		c.WS.Close()
	}
}

func (c *Conn) GetAddr() string {
	switch c.CoonType {
	case CoonTypeTCP:
		return c.TCP.RemoteAddr().String()
	case ConnTypeWS:
		return c.WS.RemoteAddr().String()
	}
	return ""
}

// HandleMessage 消息处理
func (c *Conn) HandleMessage(buf []byte) {
	var packet = new(pb.Packet)
	err := proto.Unmarshal(buf, packet)
	if err != nil {
		slog.Error("unmarshal error", "error", err, "len", len(buf))
		return
	}
	slog.Debug("HandleMessage", "packet", packet)

	// 对未登录的用户进行拦截
	if packet.Command != pb.Command_SIGN_IN && c.UserId == 0 {
		// 应该告诉用户没有登录
		return
	}

	switch packet.Command {
	case pb.Command_SIGN_IN:
		c.SignIn(packet)
	case pb.Command_SYNC:
		c.Sync(packet)
	case pb.Command_HEARTBEAT:
		c.Heartbeat(packet)
	case pb.Command_MESSAGE:
		c.MessageACK(packet)
	case pb.Command_SUBSCRIBE_ROOM:
		c.SubscribedRoom(packet)
	default:
		slog.Error("handler switch other")
	}
}

// Send 下发消息
func (c *Conn) Send(packet *pb.Packet, message proto.Message, err error) {
	packet.Data = nil
	packet.Code = 0
	packet.Message = ""

	if err != nil {
		status, _ := status.FromError(err)
		packet.Code = int32(status.Code())
		packet.Message = status.Message()
	}

	if message != nil {
		buf, err := proto.Marshal(message)
		if err != nil {
			slog.Error("proto.Marshal error", "error", err)
			return
		}
		packet.Data = buf
	}

	buf, err := proto.Marshal(packet)
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
	slog.Info("Send", "userID", c.UserId, "message", packet)
}

// SignIn 登录
func (c *Conn) SignIn(packet *pb.Packet) {
	var signIn pb.SignInInput
	err := proto.Unmarshal(packet.Data, &signIn)
	if err != nil {
		slog.Error("proto unmarshal error", "error", err)
		return
	}

	_, err = rpc.GetDeviceIntClient().ConnSignIn(md.ContextWithRequestId(context.TODO(), packet.RequestId), &logicpb.ConnSignInRequest{
		UserId:     signIn.UserId,
		DeviceId:   signIn.DeviceId,
		Token:      signIn.Token,
		ConnAddr:   config.Config.ConnectLocalAddr,
		ClientAddr: c.GetAddr(),
	})

	c.Send(packet, nil, err)
	if err != nil {
		return
	}

	c.UserId = signIn.UserId
	c.DeviceId = signIn.DeviceId
	SetConn(signIn.DeviceId, c)
}

// Sync 消息同步
func (c *Conn) Sync(packet *pb.Packet) {
	var sync pb.SyncInput
	err := proto.Unmarshal(packet.Data, &sync)
	if err != nil {
		slog.Error("proto unmarshal error", "error", err)
		return
	}
	ctx := md.ContextWithRequestId(context.TODO(), packet.RequestId)
	resp, err := rpc.GetMessageIntClient().Sync(ctx, &logicpb.SyncRequest{
		UserId:   c.UserId,
		DeviceId: c.DeviceId,
		Seq:      sync.Seq,
	})

	var message proto.Message
	if err == nil {
		message = &pb.SyncOutput{Messages: resp.Messages, HasMore: resp.HasMore}
	}
	c.Send(packet, message, err)
}

// Heartbeat 心跳
func (c *Conn) Heartbeat(packet *pb.Packet) {
	c.Send(packet, nil, nil)

	slog.Info("heartbeat", "device_id", c.DeviceId, "user_id", c.UserId)
}

// MessageACK 消息收到回执
func (c *Conn) MessageACK(packet *pb.Packet) {
	var messageACK pb.MessageACK
	err := proto.Unmarshal(packet.Data, &messageACK)
	if err != nil {
		slog.Error("proto unmarshal error", "error", err)
		return
	}

	ctx := md.ContextWithRequestId(context.TODO(), packet.RequestId)
	_, _ = rpc.GetMessageIntClient().MessageACK(ctx, &logicpb.MessageACKRequest{
		UserId:      c.UserId,
		DeviceId:    c.DeviceId,
		DeviceAck:   messageACK.DeviceAck,
		ReceiveTime: messageACK.ReceiveTime,
	})
}

// SubscribedRoom 订阅房间
func (c *Conn) SubscribedRoom(packet *pb.Packet) {
	var subscribeRoom pb.SubscribeRoomInput
	err := proto.Unmarshal(packet.Data, &subscribeRoom)
	if err != nil {
		slog.Error("proto unmarshal", "error", err)
		return
	}

	SubscribedRoom(c, subscribeRoom.RoomId)
	c.Send(packet, nil, nil)
	_, err = rpc.GetRoomIntClient().SubscribeRoom(context.TODO(), &logicpb.SubscribeRoomRequest{
		UserId:   c.UserId,
		DeviceId: c.DeviceId,
		RoomId:   subscribeRoom.RoomId,
		Seq:      subscribeRoom.Seq,
		ConnAddr: config.Config.ConnectLocalAddr,
	})
	if err != nil {
		slog.Error("SubscribedRoom error", "error", err)
	}
}

package ws

import (
	"context"
	"encoding/binary"
	"fmt"
	"gim/public/logger"
	"gim/public/pb"
	"gim/public/rpc_cli"
	"gim/public/util"
	"io"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const PreConn = -1 // 设备第二次重连时，标记设备的上一条连接
const TypeLen = 2

type WSConnContext struct {
	Conn     *websocket.Conn // websocket连接
	AppId    int64           // AppId
	DeviceId int64           // 设备id
	UserId   int64           // 用户id
}

func NewWSConnContext(conn *websocket.Conn, appId, userId, deviceId int64) *WSConnContext {
	return &WSConnContext{
		Conn:     conn,
		AppId:    appId,
		UserId:   userId,
		DeviceId: deviceId,
	}
}

// DoConn 处理连接
func (c *WSConnContext) DoConn() {
	defer util.RecoverPanic()

	for {
		err := c.Conn.SetReadDeadline(time.Now().Add(12 * time.Minute))
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		pt := int(binary.BigEndian.Uint16(data[0:2]))
		c.HandlePackage(pt, data[2:])
	}
}

// HandlePackage 处理请求发包
func (c *WSConnContext) HandlePackage(t int, bytes []byte) {
	switch pb.PackageType(t) {
	case pb.PackageType_PT_SYNC:
		c.Sync(bytes)
	case pb.PackageType_PT_HEARTBEAT:
		c.Heartbeat(bytes)
	case pb.PackageType_PT_MESSAGE:
		c.MessageACK(bytes)
	default:
		logger.Logger.Info("switch other", zap.Int("type", t))
	}

}

// Sync 离线消息同步
func (c *WSConnContext) Sync(bytes []byte) {
	var input pb.SyncInput
	err := proto.Unmarshal(bytes, &input)
	if err != nil {
		logger.Sugar.Error(err)
		c.Release()
		return
	}

	resp, err := rpc_cli.LogicIntClient.Sync(context.TODO(), &pb.SyncReq{
		AppId:    c.AppId,
		UserId:   c.UserId,
		DeviceId: c.DeviceId,
		Seq:      input.Seq,
	})

	s, _ := status.FromError(err)
	var output = pb.SyncOutput{
		Code:    int32(s.Code()),
		Message: s.Message(),
	}

	if err == nil {
		output.Messages = resp.Messages
	}

	c.Output(pb.PackageType_PT_SYNC, &output)
	if s.Code() != codes.OK {
		logger.Sugar.Error(err)
		return
	}
}

// Heartbeat 心跳
func (c *WSConnContext) Heartbeat(bytes []byte) {
	c.Output(pb.PackageType_PT_HEARTBEAT, nil)
	logger.Sugar.Infow("heartbeat", "device_id", c.DeviceId, "user_id", c.UserId)
}

// MessageACK 消息回执
func (c *WSConnContext) MessageACK(bytes []byte) {
	var input pb.MessageACK
	err := proto.Unmarshal(bytes, &input)
	if err != nil {
		logger.Sugar.Error(err)
		c.Release()
		return
	}

	_, _ = rpc_cli.LogicIntClient.MessageACK(context.TODO(), &pb.MessageACKReq{
		AppId:       c.AppId,
		UserId:      c.UserId,
		DeviceId:    c.DeviceId,
		MessageId:   input.MessageId,
		DeviceAck:   input.DeviceAck,
		ReceiveTime: input.ReceiveTime,
	})
}

// Output
func (c *WSConnContext) Output(pt pb.PackageType, message proto.Message) {
	var (
		bytes []byte
		err   error
	)

	if message != nil {
		bytes, err = proto.Marshal(message)
		if err != nil {
			logger.Sugar.Error(err)
			return
		}
	}

	writeBytes := make([]byte, len(bytes)+TypeLen)
	binary.BigEndian.PutUint16(writeBytes[0:TypeLen], uint16(pt))
	copy(writeBytes[TypeLen:], bytes)
	err = c.Conn.WriteMessage(websocket.BinaryMessage, writeBytes)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
}

// HandleReadErr 读取conn错误
func (c *WSConnContext) HandleReadErr(err error) {
	logger.Logger.Debug("read tcp error：", zap.Int64("app_id", c.AppId), zap.Int64("user_id", c.UserId),
		zap.Int64("device_id", c.DeviceId), zap.Error(err))
	str := err.Error()
	// 服务器主动关闭连接
	if strings.HasSuffix(str, "use of closed network connection") {
		return
	}

	c.Release()
	// 客户端主动关闭连接或者异常程序退出
	if err == io.EOF {
		return
	}
	// SetReadDeadline 之后，超时返回的错误
	if strings.HasSuffix(str, "i/o timeout") {
		return
	}
}

// Release 释放TCP连接
func (c *WSConnContext) Release() {
	// 从本地manager中删除tcp连接
	if c.DeviceId != PreConn {
		delete(c.DeviceId)
	}

	// 关闭tcp连接
	err := c.Conn.Close()
	if err != nil {
		logger.Sugar.Error(err)
	}

	// 通知业务服务器设备下线
	if c.DeviceId != PreConn {
		_, _ = rpc_cli.LogicIntClient.Offline(context.TODO(), &pb.OfflineReq{
			AppId:    c.AppId,
			UserId:   c.UserId,
			DeviceId: c.DeviceId,
		})
	}
}

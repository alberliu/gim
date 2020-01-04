package ws_conn

import (
	"context"
	"fmt"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc_cli"
	"gim/pkg/util"
	"io"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

const PreConn = -1 // 设备第二次重连时，标记设备的上一条连接

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

		c.HandlePackage(data)
	}
}

// HandlePackage 处理请求发包
func (c *WSConnContext) HandlePackage(bytes []byte) {
	var input pb.Input
	err := proto.Unmarshal(bytes, &input)
	if err != nil {
		logger.Sugar.Error(err)
		c.Release()
	}

	switch input.Type {
	case pb.PackageType_PT_SYNC:
		c.Sync(input)
	case pb.PackageType_PT_HEARTBEAT:
		c.Heartbeat(input)
	case pb.PackageType_PT_MESSAGE:
		c.MessageACK(input)
	default:
		logger.Logger.Info("switch other")
	}

}

// Sync 离线消息同步
func (c *WSConnContext) Sync(input pb.Input) {
	var sync pb.SyncInput
	err := proto.Unmarshal(input.Data, &sync)
	if err != nil {
		logger.Sugar.Error(err)
		c.Release()
		return
	}

	resp, err := rpc_cli.LogicIntClient.Sync(grpclib.ContextWithRequstId(context.TODO(), input.RequestId), &pb.SyncReq{
		AppId:    c.AppId,
		UserId:   c.UserId,
		DeviceId: c.DeviceId,
		Seq:      sync.Seq,
	})

	var message proto.Message
	if err == nil {
		message = &pb.SyncOutput{Messages: resp.Messages}
	}

	c.Output(pb.PackageType_PT_SYNC, input.RequestId, err, message)
}

// Heartbeat 心跳
func (c *WSConnContext) Heartbeat(input pb.Input) {
	c.Output(pb.PackageType_PT_HEARTBEAT, input.RequestId, nil, nil)
	logger.Sugar.Infow("heartbeat", "device_id", c.DeviceId, "user_id", c.UserId)
}

// MessageACK 消息回执
func (c *WSConnContext) MessageACK(input pb.Input) {
	var messageACK pb.MessageACK
	err := proto.Unmarshal(input.Data, &messageACK)
	if err != nil {
		logger.Sugar.Error(err)
		c.Release()
		return
	}

	_, _ = rpc_cli.LogicIntClient.MessageACK(grpclib.ContextWithRequstId(context.TODO(), input.RequestId), &pb.MessageACKReq{
		AppId:       c.AppId,
		UserId:      c.UserId,
		DeviceId:    c.DeviceId,
		DeviceAck:   messageACK.DeviceAck,
		ReceiveTime: messageACK.ReceiveTime,
	})
}

// Output
func (c *WSConnContext) Output(pt pb.PackageType, requestId int64, err error, message proto.Message) {
	var output = pb.Output{
		Type:      pt,
		RequestId: requestId,
	}

	if err != nil {
		status, _ := status.FromError(err)
		output.Code = int32(status.Code())
		output.Message = status.Message()
	}

	if message != nil {
		msgBytes, err := proto.Marshal(message)
		if err != nil {
			logger.Sugar.Error(err)
			return
		}
		output.Data = msgBytes
	}

	outputBytes, err := proto.Marshal(&output)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	err = c.Conn.WriteMessage(websocket.BinaryMessage, outputBytes)
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

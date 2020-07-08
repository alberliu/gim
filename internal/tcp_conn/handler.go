package tcp_conn

import (
	"context"
	"gim/config"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc_cli"

	"github.com/alberliu/gn"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

type ConnData struct {
	AppId    int64 // AppId
	DeviceId int64 // 设备id
	UserId   int64 // 用户id
}

const PreConn = -1 // 设备第二次重连时，标记设备的上一条连接

type handler struct{}

var Handler = new(handler)

func (*handler) OnConnect(c *gn.Conn) {
	logger.Logger.Debug("connect:", zap.Int("fd", c.GetFd()), zap.String("addr", c.GetAddr()))
}

func (h *handler) OnMessage(c *gn.Conn, bytes []byte) {
	var input pb.Input
	err := proto.Unmarshal(bytes, &input)
	if err != nil {
		logger.Logger.Error("unmarshal error", zap.Error(err))
		return
	}

	// 对未登录的用户进行拦截
	if input.Type != pb.PackageType_PT_SIGN_IN && c.GetData() == nil {
		// 应该告诉用户没有登录
		return
	}

	switch input.Type {
	case pb.PackageType_PT_SIGN_IN:
		h.SignIn(c, input)
	case pb.PackageType_PT_SYNC:
		h.Sync(c, input)
	case pb.PackageType_PT_HEARTBEAT:
		h.Heartbeat(c, input)
	case pb.PackageType_PT_MESSAGE:
		h.MessageACK(c, input)
	default:
		logger.Logger.Error("handler switch other")
	}
	return
}

func (*handler) OnClose(c *gn.Conn, err error) {
	logger.Logger.Debug("close", zap.Any("data", c.GetData()), zap.Error(err))
	data := c.GetData().(ConnData)
	_, _ = rpc_cli.LogicIntClient.Offline(context.TODO(), &pb.OfflineReq{
		AppId:    data.AppId,
		UserId:   data.UserId,
		DeviceId: data.DeviceId,
	})
}

// SignIn 登录
func (*handler) SignIn(c *gn.Conn, input pb.Input) {
	var signIn pb.SignInInput
	err := proto.Unmarshal(input.Data, &signIn)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	_, err = rpc_cli.LogicIntClient.SignIn(grpclib.ContextWithRequstId(context.TODO(), input.RequestId), &pb.SignInReq{
		AppId:    signIn.AppId,
		UserId:   signIn.UserId,
		DeviceId: signIn.DeviceId,
		Token:    signIn.Token,
		ConnAddr: config.TCPConnConf.LocalAddr,
		ConnFd:   int64(c.GetFd()),
	})

	send(c, pb.PackageType_PT_SIGN_IN, input.RequestId, err, nil)

	data := ConnData{
		AppId:    signIn.AppId,
		DeviceId: signIn.UserId,
		UserId:   signIn.DeviceId,
	}
	c.SetData(data)
}

// Sync 消息同步
func (*handler) Sync(c *gn.Conn, input pb.Input) {
	var sync pb.SyncInput
	err := proto.Unmarshal(input.Data, &sync)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	data := c.GetData().(ConnData)
	resp, err := rpc_cli.LogicIntClient.Sync(grpclib.ContextWithRequstId(context.TODO(), input.RequestId), &pb.SyncReq{
		AppId:    data.AppId,
		UserId:   data.UserId,
		DeviceId: data.DeviceId,
		Seq:      sync.Seq,
	})

	var message proto.Message
	if err == nil {
		message = &pb.SyncOutput{Messages: resp.Messages}
	}
	send(c, pb.PackageType_PT_SYNC, input.RequestId, err, message)
}

// Heartbeat 心跳
func (*handler) Heartbeat(c *gn.Conn, input pb.Input) {
	data := c.GetData().(ConnData)
	send(c, pb.PackageType_PT_HEARTBEAT, input.RequestId, nil, nil)
	logger.Sugar.Infow("heartbeat", "device_id", data.DeviceId, "user_id", data.UserId)
}

// MessageACK 消息收到回执
func (*handler) MessageACK(c *gn.Conn, input pb.Input) {
	var messageACK pb.MessageACK
	err := proto.Unmarshal(input.Data, &messageACK)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	data := c.GetData().(ConnData)
	_, _ = rpc_cli.LogicIntClient.MessageACK(grpclib.ContextWithRequstId(context.TODO(), input.RequestId), &pb.MessageACKReq{
		AppId:       data.AppId,
		UserId:      data.UserId,
		DeviceId:    data.DeviceId,
		DeviceAck:   messageACK.DeviceAck,
		ReceiveTime: messageACK.ReceiveTime,
	})
}

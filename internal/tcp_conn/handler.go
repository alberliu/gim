package tcp_conn

import (
	"context"
	"gim/config"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc"

	"github.com/alberliu/gn"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

type ConnData struct {
	UserId   int64
	DeviceId int64
}

type handler struct{}

var Handler = new(handler)

func (*handler) OnConnect(c *gn.Conn) {
	logger.Logger.Debug("connect:", zap.Int32("fd", c.GetFd()), zap.String("addr", c.GetAddr()))
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
	logger.Logger.Debug("close", zap.String("addr", c.GetAddr()), zap.Int32("fd", c.GetFd()), zap.Any("data", c.GetData()), zap.Error(err))
	if data, ok := c.GetData().(ConnData); ok {
		_, _ = rpc.LogicIntClient.Offline(context.TODO(), &pb.OfflineReq{
			UserId:     data.UserId,
			DeviceId:   data.DeviceId,
			ClientAddr: c.GetAddr(),
		})
	}
}

func (h *handler) Send(c *gn.Conn, pt pb.PackageType, requestId int64, err error, message proto.Message) {
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

	err = encoder.EncodeToFD(c.GetFd(), outputBytes)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
}

// SignIn 登录
func (h *handler) SignIn(c *gn.Conn, input pb.Input) {
	var signIn pb.SignInInput
	err := proto.Unmarshal(input.Data, &signIn)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	_, err = rpc.LogicIntClient.ConnSignIn(grpclib.ContextWithRequstId(context.TODO(), input.RequestId), &pb.ConnSignInReq{
		UserId:     signIn.UserId,
		DeviceId:   signIn.DeviceId,
		Token:      signIn.Token,
		ConnAddr:   config.TCPConn.LocalAddr,
		ConnFd:     int64(c.GetFd()),
		ClientAddr: c.GetAddr(),
	})

	h.Send(c, pb.PackageType_PT_SIGN_IN, input.RequestId, err, nil)
	if err != nil {
		return
	}

	data := ConnData{
		UserId:   signIn.UserId,
		DeviceId: signIn.DeviceId,
	}
	c.SetData(data)
}

// Sync 消息同步
func (h *handler) Sync(c *gn.Conn, input pb.Input) {
	var sync pb.SyncInput
	err := proto.Unmarshal(input.Data, &sync)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	data := c.GetData().(ConnData)
	resp, err := rpc.LogicIntClient.Sync(grpclib.ContextWithRequstId(context.TODO(), input.RequestId), &pb.SyncReq{
		UserId:   data.UserId,
		DeviceId: data.DeviceId,
		Seq:      sync.Seq,
	})

	var message proto.Message
	if err == nil {
		message = &pb.SyncOutput{Messages: resp.Messages, HasMore: resp.HasMore}
	}
	h.Send(c, pb.PackageType_PT_SYNC, input.RequestId, err, message)
}

// Heartbeat 心跳
func (h *handler) Heartbeat(c *gn.Conn, input pb.Input) {
	h.Send(c, pb.PackageType_PT_HEARTBEAT, input.RequestId, nil, nil)
	data := c.GetData().(ConnData)
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
	_, _ = rpc.LogicIntClient.MessageACK(grpclib.ContextWithRequstId(context.TODO(), input.RequestId), &pb.MessageACKReq{
		UserId:      data.UserId,
		DeviceId:    data.DeviceId,
		DeviceAck:   messageACK.DeviceAck,
		ReceiveTime: messageACK.ReceiveTime,
	})
}

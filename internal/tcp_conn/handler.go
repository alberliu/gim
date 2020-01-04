package tcp_conn

import (
	"context"
	"gim/config"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc_cli"

	"go.uber.org/zap"

	"github.com/golang/protobuf/proto"
)

const PreConn = -1 // 设备第二次重连时，标记设备的上一条连接

type handler struct{}

var Handler = new(handler)

// Handler 处理客户端的上行包
func (h *handler) Handler(ctx *ConnContext, bytes []byte) {
	var input pb.Input
	err := proto.Unmarshal(bytes, &input)
	if err != nil {
		logger.Logger.Error("unmarshal error", zap.Error(err))
		ctx.Release()
		return
	}

	// 对未登录的用户进行拦截
	if input.Type != pb.PackageType_PT_SIGN_IN && ctx.IsSignIn == false {
		// 应该告诉用户没有登录
		ctx.Release()
		return
	}

	switch input.Type {
	case pb.PackageType_PT_SIGN_IN:
		h.SignIn(ctx, input)
	case pb.PackageType_PT_SYNC:
		h.Sync(ctx, input)
	case pb.PackageType_PT_HEARTBEAT:
		h.Heartbeat(ctx, input)
	case pb.PackageType_PT_MESSAGE:
		h.MessageACK(ctx, input)
	default:
		logger.Logger.Error("handler switch other")
	}
	return
}

// SignIn 登录
func (*handler) SignIn(ctx *ConnContext, input pb.Input) {
	var signIn pb.SignInInput
	err := proto.Unmarshal(input.Data, &signIn)
	if err != nil {
		logger.Sugar.Error(err)
		ctx.Release()
		return
	}

	_, err = rpc_cli.LogicIntClient.SignIn(grpclib.ContextWithRequstId(context.TODO(), input.RequestId), &pb.SignInReq{
		AppId:    signIn.AppId,
		UserId:   signIn.UserId,
		DeviceId: signIn.DeviceId,
		Token:    signIn.Token,
		ConnAddr: config.ConnConf.LocalAddr,
	})

	ctx.Output(pb.PackageType_PT_SIGN_IN, input.RequestId, err, nil)
	if err != nil {
		ctx.Release()
		return
	}

	ctx.AppId = signIn.AppId
	ctx.UserId = signIn.UserId
	ctx.DeviceId = signIn.DeviceId
	ctx.IsSignIn = true

	// 断开这个设备之前的连接
	preCtx := load(ctx.DeviceId)
	if preCtx != nil {
		preCtx.DeviceId = PreConn
	}

	store(ctx.DeviceId, ctx)
}

// Sync 消息同步
func (*handler) Sync(ctx *ConnContext, input pb.Input) {
	var sync pb.SyncInput
	err := proto.Unmarshal(input.Data, &sync)
	if err != nil {
		logger.Sugar.Error(err)
		ctx.Release()
		return
	}

	resp, err := rpc_cli.LogicIntClient.Sync(grpclib.ContextWithRequstId(context.TODO(), input.RequestId), &pb.SyncReq{
		AppId:    ctx.AppId,
		UserId:   ctx.UserId,
		DeviceId: ctx.DeviceId,
		Seq:      sync.Seq,
	})

	var message proto.Message
	if err == nil {
		message = &pb.SyncOutput{Messages: resp.Messages}
	}
	ctx.Output(pb.PackageType_PT_SYNC, input.RequestId, err, message)
}

// Heartbeat 心跳
func (*handler) Heartbeat(ctx *ConnContext, input pb.Input) {
	ctx.Output(pb.PackageType_PT_HEARTBEAT, input.RequestId, nil, nil)
	logger.Sugar.Infow("heartbeat", "device_id", ctx.DeviceId, "user_id", ctx.UserId)
}

// MessageACK 消息收到回执
func (*handler) MessageACK(ctx *ConnContext, input pb.Input) {
	var messageACK pb.MessageACK
	err := proto.Unmarshal(input.Data, &messageACK)
	if err != nil {
		logger.Sugar.Error(err)
		ctx.Release()
		return
	}

	_, _ = rpc_cli.LogicIntClient.MessageACK(grpclib.ContextWithRequstId(context.TODO(), input.RequestId), &pb.MessageACKReq{
		AppId:       ctx.AppId,
		UserId:      ctx.UserId,
		DeviceId:    ctx.DeviceId,
		DeviceAck:   messageACK.DeviceAck,
		ReceiveTime: messageACK.ReceiveTime,
	})
}

// Offline 设备离线
func (*handler) Offline(ctx *ConnContext) {
	_, _ = rpc_cli.LogicIntClient.Offline(context.TODO(), &pb.OfflineReq{
		AppId:    ctx.AppId,
		UserId:   ctx.UserId,
		DeviceId: ctx.DeviceId,
	})
}

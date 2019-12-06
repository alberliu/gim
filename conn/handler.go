package conn

import (
	"context"
	"gim/conf"
	"gim/public/logger"
	"gim/public/pb"
	"gim/public/rpc_cli"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const PreConn = -1 // 设备第二次重连时，标记设备的上一条连接

type handler struct{}

var Handler = new(handler)

// Handler 处理客户端的上行包
func (h *handler) Handler(ctx *ConnContext, pack *Package) {
	// 对未登录的用户进行拦截
	if pack.Code != int(pb.PackageType_PT_SIGN_IN) && ctx.IsSignIn == false {
		// 应该告诉用户没有登录
		ctx.Release()
		return
	}

	switch pb.PackageType(pack.Code) {
	case pb.PackageType_PT_SIGN_IN:
		h.SignIn(ctx, pack.Content)
	case pb.PackageType_PT_SYNC:
		h.Sync(ctx, pack.Content)
	case pb.PackageType_PT_HEARTBEAT:
		h.Heartbeat(ctx, pack.Content)
	case pb.PackageType_PT_MESSAGE:
		h.MessageACK(ctx, pack.Content)
	}
	return
}

// SignIn 登录
func (*handler) SignIn(ctx *ConnContext, bytes []byte) {
	var input pb.SignInInput
	err := proto.Unmarshal(bytes, &input)
	if err != nil {
		logger.Sugar.Error(err)
		ctx.Release()
		return
	}

	_, err = rpc_cli.LogicIntClient.SignIn(context.TODO(), &pb.SignInReq{
		AppId:    input.AppId,
		UserId:   input.UserId,
		DeviceId: input.DeviceId,
		Token:    input.Token,
		ConnAddr: conf.LocalAddr,
	})

	s, _ := status.FromError(err)
	ctx.Output(pb.PackageType_PT_SIGN_IN, &pb.SignInOutput{Code: int32(s.Code()), Message: s.Message()})
	if s.Code() != codes.OK {
		ctx.Release()
		return
	}

	ctx.AppId = input.AppId
	ctx.UserId = input.UserId
	ctx.DeviceId = input.DeviceId
	ctx.IsSignIn = true

	// 断开这个设备之前的连接
	preCtx := load(ctx.DeviceId)
	if preCtx != nil {
		preCtx.DeviceId = PreConn
	}

	store(ctx.DeviceId, ctx)
}

// Sync 消息同步
func (*handler) Sync(ctx *ConnContext, bytes []byte) {
	var input pb.SyncInput
	err := proto.Unmarshal(bytes, &input)
	if err != nil {
		logger.Sugar.Error(err)
		ctx.Release()
		return
	}

	resp, err := rpc_cli.LogicIntClient.Sync(context.TODO(), &pb.SyncReq{
		AppId:    ctx.AppId,
		UserId:   ctx.UserId,
		DeviceId: ctx.DeviceId,
		Seq:      input.Seq,
	})

	s, _ := status.FromError(err)
	var output = pb.SyncOutput{
		Code:     int32(s.Code()),
		Message:  s.Message(),
		Messages: resp.Messages,
	}

	ctx.Output(pb.PackageType_PT_SYNC, &output)
	if s.Code() != codes.OK {
		logger.Sugar.Error(err)
		return
	}
}

// Heartbeat 心跳
func (*handler) Heartbeat(ctx *ConnContext, bytes []byte) {
	ctx.Output(pb.PackageType_PT_HEARTBEAT, nil)
	logger.Sugar.Infow("heartbeat", "device_id", ctx.DeviceId, "user_id", ctx.UserId)
}

// MessageACK 消息收到回执
func (*handler) MessageACK(ctx *ConnContext, bytes []byte) {
	var input pb.MessageACK
	err := proto.Unmarshal(bytes, &input)
	if err != nil {
		logger.Sugar.Error(err)
		ctx.Release()
		return
	}

	_, _ = rpc_cli.LogicIntClient.MessageACK(context.TODO(), &pb.MessageACKReq{
		AppId:       ctx.AppId,
		UserId:      ctx.UserId,
		DeviceId:    ctx.DeviceId,
		MessageId:   input.MessageId,
		DeviceAck:   input.DeviceAck,
		ReceiveTime: input.ReceiveTime,
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

package rpc

import (
	"context"
	"gim/logic/model"
	"gim/logic/service"
	"gim/public/imctx"
	"gim/public/logger"
	"gim/public/pb"

	"go.uber.org/zap"
)

func Context() *imctx.Context {
	return imctx.NewContext()
}

type LogicIntServer struct{}

// SignIn 设备登录
func (*LogicIntServer) SignIn(ctx context.Context, req *pb.SignInReq) (*pb.SignInResp, error) {
	logger.Logger.Debug("device sign_in req", zap.Int64("app_id", req.AppId),
		zap.Int64("user_id", req.UserId), zap.Int64("device_id", req.DeviceId))

	err := service.AuthService.SignIn(Context(), req.AppId, req.UserId, req.DeviceId, req.Token, req.ConnAddr)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	logger.Logger.Debug("device sign_in resp", zap.Int64("app_id", req.AppId),
		zap.Int64("user_id", req.UserId), zap.Int64("device_id", req.DeviceId),
		zap.Error(err))
	return &pb.SignInResp{}, nil
}

// Sync 设备同步消息
func (*LogicIntServer) Sync(ctx context.Context, req *pb.SyncReq) (*pb.SyncResp, error) {
	logger.Logger.Debug("sync req", zap.Int64("app_id", req.AppId),
		zap.Int64("user_id", req.UserId), zap.Int64("device_id", req.DeviceId),
		zap.Int64("seq", req.Seq))

	messages, err := service.MessageService.ListByUserIdAndSeq(Context(), req.AppId, req.UserId, req.Seq)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	logger.Logger.Debug("sync resp", zap.Int64("app_id", req.AppId),
		zap.Int64("user_id", req.UserId), zap.Int64("device_id", req.DeviceId),
		zap.Int64s("seqs", service.MessageService.GetMessageSeqs(messages)), zap.Error(err))
	return &pb.SyncResp{Messages: model.MessagesToPB(messages)}, nil
}

// MessageACK 设备收到消息ack
func (*LogicIntServer) MessageACK(ctx context.Context, req *pb.MessageACKReq) (*pb.MessageACKResp, error) {
	err := service.DeviceAckService.Update(Context(), req.DeviceId, req.DeviceAck)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	logger.Logger.Debug("message ack", zap.String("message_id", req.MessageId),
		zap.Int64("app_id", req.AppId), zap.Int64("user_id", req.UserId),
		zap.Int64("device_id", req.DeviceId), zap.Int64("ack", req.DeviceAck))
	return &pb.MessageACKResp{}, nil
}

// Offline 设备离线
func (*LogicIntServer) Offline(ctx context.Context, req *pb.OfflineReq) (*pb.OfflineResp, error) {
	err := service.DeviceService.Offline(Context(), req.AppId, req.UserId, req.DeviceId)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return &pb.OfflineResp{}, nil
}

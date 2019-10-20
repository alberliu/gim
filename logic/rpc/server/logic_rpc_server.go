package server

import (
	"gim/logic/model"
	"gim/logic/service"
	"gim/public/imctx"
	"gim/public/imerror"
	"gim/public/logger"
	"gim/public/pb"
	"gim/public/transfer"

	"go.uber.org/zap"

	"github.com/golang/protobuf/proto"
)

func Context() *imctx.Context {
	return imctx.NewContext()
}

type LogicRPCServer struct{}

// SignIn 设备登录
func (s *LogicRPCServer) SignIn(req transfer.SignInReq, resp *transfer.SignInResp) error {
	signInReq := pb.SignInReq{}
	err := proto.Unmarshal(req.Bytes, &signInReq)
	if err != nil {
		logger.Sugar.Error(err)
		*resp = *transfer.ErrorToSignInResp(err, 0, 0, 0)
		return nil
	}

	logger.Logger.Debug("device sign_in req", zap.Int64("app_id", signInReq.AppId),
		zap.Int64("user_id", signInReq.UserId), zap.Int64("device_id", signInReq.DeviceId))

	err = service.AuthService.SignIn(Context(), signInReq.AppId, signInReq.UserId, signInReq.DeviceId, signInReq.Token, req.ConnectIP)
	if err != nil {
		logger.Sugar.Error(err)
	}

	*resp = *transfer.ErrorToSignInResp(err, signInReq.AppId, signInReq.UserId, signInReq.DeviceId)

	logger.Logger.Debug("device sign_in resp", zap.Int64("app_id", signInReq.AppId),
		zap.Int64("user_id", signInReq.UserId), zap.Int64("device_id", signInReq.DeviceId),
		zap.Error(err))
	return nil
}

// Sync 设备同步消息
func (s *LogicRPCServer) Sync(req transfer.SyncReq, resp *transfer.SyncResp) error {
	if !req.IsSignIn {
		*resp = *transfer.ErrorToSyncResp(imerror.ErrUnauthorized, nil)
		return nil
	}

	syncReq := pb.SyncReq{}
	err := proto.Unmarshal(req.Bytes, &syncReq)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	logger.Logger.Debug("sync req", zap.Int64("app_id", req.AppId),
		zap.Int64("user_id", req.UserId), zap.Int64("device_id", req.DeviceId),
		zap.Int64("seq", syncReq.Seq))

	messages, err := service.MessageService.ListByUserIdAndSeq(Context(), req.AppId, req.UserId, syncReq.Seq)
	if err != nil {
		logger.Sugar.Error(err)
		*resp = *transfer.ErrorToSyncResp(imerror.ErrUnauthorized, nil)
		return nil
	}
	logger.Logger.Debug("sync resp", zap.Int64("app_id", req.AppId),
		zap.Int64("user_id", req.UserId), zap.Int64("device_id", req.DeviceId),
		zap.Int64s("seqs", service.MessageService.GetMessageSeqs(messages)), zap.Error(err))
	*resp = *transfer.ErrorToSyncResp(nil, model.MessagesToPB(messages))
	return nil
}

// MessageACK 设备收到消息ack
func (s *LogicRPCServer) MessageACK(req transfer.MessageAckReq, resp *transfer.MessageAckResp) error {
	messageACK := pb.MessageACK{}
	err := proto.Unmarshal(req.Bytes, &messageACK)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	err = service.DeviceAckService.Update(Context(), req.DeviceId, messageACK.DeviceAck)
	if err != nil {
		logger.Sugar.Error(err)
	}

	logger.Logger.Debug("message ack", zap.String("message_id", messageACK.MessageId),
		zap.Int64("app_id", req.AppId), zap.Int64("user_id", req.UserId),
		zap.Int64("device_id", req.DeviceId), zap.Int64("ack", messageACK.DeviceAck))
	return nil
}

// Offline 设备离线
func (s *LogicRPCServer) Offline(req transfer.OfflineReq, resp *transfer.OfflineResp) error {
	err := service.DeviceService.Offline(Context(), req.AppId, req.UserId, req.DeviceId)
	if err != nil {
		logger.Sugar.Error(err)
	}
	return nil
}

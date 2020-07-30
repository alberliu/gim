package api

import (
	"context"
	"gim/internal/logic/model"
	"gim/internal/logic/service"
	"gim/pkg/gerrors"
	"gim/pkg/logger"
	"gim/pkg/pb"
)

type LogicIntServer struct{}

// SignIn 设备登录
func (*LogicIntServer) ConnSignIn(ctx context.Context, req *pb.ConnSignInReq) (*pb.ConnSignInResp, error) {
	return &pb.ConnSignInResp{}, service.AuthService.SignIn(ctx, req.UserId, req.DeviceId, req.Token, req.ConnAddr, req.ConnFd)
}

// Sync 设备同步消息
func (*LogicIntServer) Sync(ctx context.Context, req *pb.SyncReq) (*pb.SyncResp, error) {
	messages, err := service.MessageService.ListByUserIdAndSeq(ctx, req.UserId, req.Seq)
	if err != nil {
		return nil, err
	}
	return &pb.SyncResp{Messages: model.MessagesToPB(messages)}, nil
}

// MessageACK 设备收到消息ack
func (*LogicIntServer) MessageACK(ctx context.Context, req *pb.MessageACKReq) (*pb.MessageACKResp, error) {
	return &pb.MessageACKResp{}, service.DeviceAckService.Update(ctx, req.DeviceId, req.DeviceAck)
}

// Offline 设备离线
func (*LogicIntServer) Offline(ctx context.Context, req *pb.OfflineReq) (*pb.OfflineResp, error) {
	return &pb.OfflineResp{}, service.DeviceService.Offline(ctx, req.UserId, req.DeviceId)
}

// SendMessage 发送消息
func (*LogicIntServer) SendMessage(ctx context.Context, req *pb.SendMessageReq) (*pb.SendMessageResp, error) {
	sender := model.Sender{
		SenderType: pb.SenderType_ST_BUSINESS,
		SenderId:   0,
		DeviceId:   0,
	}
	seq, err := service.MessageService.Send(ctx, sender, *req)
	if err != nil {
		return nil, err
	}
	return &pb.SendMessageResp{Seq: seq}, nil
}

// GetDevice 获取设备信息
func (*LogicIntServer) GetDevice(ctx context.Context, req *pb.GetDeviceReq) (*pb.GetDeviceResp, error) {
	device, err := service.DeviceService.Get(ctx, req.DeviceId)
	if err != nil {
		return nil, err
	}

	if device == nil {
		return nil, gerrors.ErrDeviceNotExist
	}

	return &pb.GetDeviceResp{
		Device: &pb.Device{
			DeviceId:      device.Id,
			UserId:        device.UserId,
			Type:          device.Type,
			Brand:         device.Brand,
			Model:         device.Model,
			SystemVersion: device.SystemVersion,
			SDKVersion:    device.SDKVersion,
			Status:        device.Status,
			CreateTime:    device.CreateTime.Unix(),
			UpdateTime:    device.UpdateTime.Unix(),
		},
	}, nil
}

// ServerStop 服务停止
func (s *LogicIntServer) ServerStop(ctx context.Context, in *pb.ServerStopReq) (*pb.ServerStopResp, error) {
	go func() {
		err := service.DeviceService.ServerStop(ctx, in.ConnAddr)
		if err != nil {
			logger.Sugar.Error(err)
		}
	}()
	return &pb.ServerStopResp{}, nil
}

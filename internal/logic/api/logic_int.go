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

// ConnSignIn 设备登录
func (*LogicIntServer) ConnSignIn(ctx context.Context, req *pb.ConnSignInReq) (*pb.Empty, error) {
	return &pb.Empty{},
		service.AuthService.SignIn(ctx, req.UserId, req.DeviceId, req.Token, req.ConnAddr, req.ClientAddr)
}

// Sync 设备同步消息
func (*LogicIntServer) Sync(ctx context.Context, req *pb.SyncReq) (*pb.SyncResp, error) {
	return service.MessageService.Sync(ctx, req.UserId, req.Seq)
}

// MessageACK 设备收到消息ack
func (*LogicIntServer) MessageACK(ctx context.Context, req *pb.MessageACKReq) (*pb.Empty, error) {
	return &pb.Empty{}, service.DeviceAckService.Update(ctx, req.UserId, req.DeviceId, req.DeviceAck)
}

// Offline 设备离线
func (*LogicIntServer) Offline(ctx context.Context, req *pb.OfflineReq) (*pb.Empty, error) {
	return &pb.Empty{}, service.DeviceService.Offline(ctx, req.UserId, req.DeviceId, req.ClientAddr)
}

func (s *LogicIntServer) SubscribeRoom(ctx context.Context, req *pb.SubscribeRoomReq) (*pb.Empty, error) {
	return &pb.Empty{}, service.RoomService.SubscribeRoom(ctx, req)
}

// SendMessage 发送消息
func (*LogicIntServer) SendMessage(ctx context.Context, req *pb.SendMessageReq) (*pb.SendMessageResp, error) {
	sender := model.Sender{
		SenderType: pb.SenderType_ST_BUSINESS,
		SenderId:   0,
		DeviceId:   0,
	}
	seq, err := service.MessageService.Send(ctx, sender, req)
	if err != nil {
		return nil, err
	}
	return &pb.SendMessageResp{Seq: seq}, nil
}

// PushRoom 推送房间
func (s *LogicIntServer) PushRoom(ctx context.Context, req *pb.PushRoomReq) (*pb.Empty, error) {
	return &pb.Empty{}, service.RoomService.Push(ctx, model.Sender{
		SenderType: pb.SenderType_ST_BUSINESS,
	}, req)
}

// PushAll 全服推送
func (s *LogicIntServer) PushAll(ctx context.Context, req *pb.PushAllReq) (*pb.Empty, error) {
	return &pb.Empty{}, service.PushService.PushAll(ctx, req)
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
func (s *LogicIntServer) ServerStop(ctx context.Context, in *pb.ServerStopReq) (*pb.Empty, error) {
	go func() {
		err := service.DeviceService.ServerStop(ctx, in.ConnAddr)
		if err != nil {
			logger.Sugar.Error(err)
		}
	}()
	return &pb.Empty{}, nil
}

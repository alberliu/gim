package api

import (
	"context"
	"gim/internal/logic/domain/device"
	"gim/internal/logic/domain/message"
	"gim/internal/logic/domain/room"
	"gim/internal/logic/proxy"
	"gim/pkg/logger"
	"gim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type LogicIntServer struct {
	pb.UnsafeLogicIntServer
}

// ConnSignIn 设备登录
func (*LogicIntServer) ConnSignIn(ctx context.Context, req *pb.ConnSignInReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{},
		device.App.SignIn(ctx, req.UserId, req.DeviceId, req.Token, req.ConnAddr, req.ClientAddr)
}

// Sync 设备同步消息
func (*LogicIntServer) Sync(ctx context.Context, req *pb.SyncReq) (*pb.SyncResp, error) {
	return message.App.Sync(ctx, req.UserId, req.Seq)
}

// MessageACK 设备收到消息ack
func (*LogicIntServer) MessageACK(ctx context.Context, req *pb.MessageACKReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, message.App.MessageAck(ctx, req.UserId, req.DeviceId, req.DeviceAck)
}

// Offline 设备离线
func (*LogicIntServer) Offline(ctx context.Context, req *pb.OfflineReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, device.App.Offline(ctx, req.DeviceId, req.ClientAddr)
}

func (s *LogicIntServer) SubscribeRoom(ctx context.Context, req *pb.SubscribeRoomReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, room.App.SubscribeRoom(ctx, req)
}

// Push 推送
func (*LogicIntServer) Push(ctx context.Context, req *pb.PushReq) (*pb.PushResp, error) {
	seq, err := proxy.PushToUserBytes(ctx, req.UserId, req.Code, req.Content, req.IsPersist)
	if err != nil {
		return nil, err
	}
	return &pb.PushResp{Seq: seq}, nil
}

// PushRoom 推送房间
func (s *LogicIntServer) PushRoom(ctx context.Context, req *pb.PushRoomReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, room.App.Push(ctx, req)
}

// PushAll 全服推送
func (s *LogicIntServer) PushAll(ctx context.Context, req *pb.PushAllReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, message.App.PushAll(ctx, req)
}

// GetDevice 获取设备信息
func (*LogicIntServer) GetDevice(ctx context.Context, req *pb.GetDeviceReq) (*pb.GetDeviceResp, error) {
	device, err := device.App.GetDevice(ctx, req.DeviceId)
	return &pb.GetDeviceResp{Device: device}, err
}

// ServerStop 服务停止
func (s *LogicIntServer) ServerStop(ctx context.Context, in *pb.ServerStopReq) (*emptypb.Empty, error) {
	go func() {
		err := device.App.ServerStop(ctx, in.ConnAddr)
		if err != nil {
			logger.Sugar.Error(err)
		}
	}()
	return &emptypb.Empty{}, nil
}

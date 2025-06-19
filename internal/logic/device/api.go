package device

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "gim/pkg/protocol/pb/logicpb"
)

type DeviceExtService struct {
	pb.UnsafeDeviceExtServiceServer
}

// RegisterDevice 注册设备
func (*DeviceExtService) RegisterDevice(ctx context.Context, request *pb.RegisterDeviceRequest) (*pb.RegisterDeviceReply, error) {
	deviceID, err := App.Register(ctx, request)
	return &pb.RegisterDeviceReply{DeviceId: deviceID}, err
}

type DeviceIntService struct {
	pb.UnsafeDeviceIntServiceServer
}

// ConnSignIn 设备登录
func (*DeviceIntService) ConnSignIn(ctx context.Context, request *pb.ConnSignInRequest) (*emptypb.Empty, error) {
	err := App.SignIn(ctx, request.UserId, request.DeviceId, request.Token, request.ConnAddr, request.ClientAddr)
	return &emptypb.Empty{}, err
}

// Offline 设备离线
func (*DeviceIntService) Offline(ctx context.Context, request *pb.OfflineRequest) (*emptypb.Empty, error) {
	err := App.Offline(ctx, request.DeviceId, request.ClientAddr)
	return &emptypb.Empty{}, err
}

// GetDevice 获取设备信息
func (*DeviceIntService) GetDevice(ctx context.Context, request *pb.GetDeviceRequest) (*pb.GetDeviceReply, error) {
	device, err := App.GetDevice(ctx, request.DeviceId)
	return &pb.GetDeviceReply{Device: device}, err
}

// ServerStop 服务停止
func (s *DeviceIntService) ServerStop(ctx context.Context, request *pb.ServerStopRequest) (*emptypb.Empty, error) {
	go func() {
		err := App.ServerStop(ctx, request.ConnAddr)
		if err != nil {
			slog.Error("ServerStop error", "error", err)
		}
	}()
	return &emptypb.Empty{}, nil
}

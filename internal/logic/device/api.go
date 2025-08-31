package device

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "gim/pkg/protocol/pb/logicpb"
)

type DeviceIntService struct {
	pb.UnsafeDeviceIntServiceServer
}

// SignIn 设备登录
func (*DeviceIntService) SignIn(ctx context.Context, request *pb.SignInRequest) (*pb.SignInReply, error) {
	err := App.SignIn(ctx, request)
	return &pb.SignInReply{}, err
}

func (s *DeviceIntService) Heartbeat(ctx context.Context, request *pb.HeartbeatRequest) (*emptypb.Empty, error) {
	err := App.Heartbeat(ctx, request.UserId, request.DeviceId)
	return &emptypb.Empty{}, err
}

// Offline 设备离线
func (*DeviceIntService) Offline(ctx context.Context, request *pb.OfflineRequest) (*emptypb.Empty, error) {
	err := App.Offline(ctx, request.DeviceId, request.ClientAddr)
	return &emptypb.Empty{}, err
}

// Save 保存
func (*DeviceIntService) Save(ctx context.Context, request *pb.DeviceSaveRequest) (*pb.DeviceSaveReply, error) {
	deviceID, err := App.Save(ctx, request.Device)
	return &pb.DeviceSaveReply{DeviceId: deviceID}, err
}

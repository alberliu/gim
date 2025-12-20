package app

import (
	"context"

	"gim/internal/logic/device/domain"
	"gim/internal/logic/device/repo"
	"gim/pkg/md"
	"gim/pkg/protocol/pb/businesspb"
	pb "gim/pkg/protocol/pb/logicpb"
	"gim/pkg/rpc"
)

var DeviceApp = new(deviceApp)

type deviceApp struct{}

// SignIn 登录
func (*deviceApp) SignIn(ctx context.Context, request *pb.SignInRequest) error {
	_, err := rpc.GetUserIntClient().Auth(ctx, &businesspb.AuthRequest{
		UserId:   request.UserId,
		DeviceId: request.DeviceId,
		Token:    request.Token,
	})
	if err != nil {
		return err
	}

	device, err := repo.DeviceRepo.Get(ctx, request.DeviceId)
	if err != nil {
		return err
	}
	device.UserID = request.UserId
	device.ConnectIP = md.GetClientIP(ctx)
	device.ClientAddr = request.ClientAddr
	err = repo.DeviceRepo.Save(ctx, device)
	if err != nil {
		return err
	}

	return repo.DeviceRepo.SetOnline(ctx, request.DeviceId)
}

// Heartbeat 设备离线
func (*deviceApp) Heartbeat(ctx context.Context, userID, deviceID uint64) error {
	return repo.DeviceRepo.SetOnline(ctx, deviceID)
}

// Offline 设备离线
func (*deviceApp) Offline(ctx context.Context, deviceID uint64, clientAddr string) error {
	return repo.DeviceRepo.SetOffline(ctx, deviceID)
}

// ListByUserID 获取用户所有在线设备
func (*deviceApp) ListByUserID(ctx context.Context, userID uint64) ([]domain.Device, error) {
	return repo.DeviceRepo.ListByUserID(ctx, userID)
}

// Save 获取设备信息
func (*deviceApp) Save(ctx context.Context, pbdevice *pb.Device) (uint64, error) {
	device := &domain.Device{
		ID:            pbdevice.Id,
		Type:          pbdevice.Type,
		Brand:         pbdevice.Brand,
		Model:         pbdevice.Model,
		SystemVersion: pbdevice.SystemVersion,
		SDKVersion:    pbdevice.SdkVersion,
		BrandPushID:   pbdevice.BranchPushId,
	}

	err := repo.DeviceRepo.Save(ctx, device)
	return device.ID, err
}

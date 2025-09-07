package app

import (
	"context"

	"gim/internal/logic/device/domain"
	"gim/internal/logic/device/repo"
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

	device, err := repo.DeviceRepo.Get(request.DeviceId)
	if err != nil {
		return err
	}
	device.UserID = request.UserId
	device.ConnectAddr = request.ConnectAddr
	device.ClientAddr = request.ClientAddr
	err = repo.DeviceRepo.Save(device)
	if err != nil {
		return err
	}

	return repo.DeviceRepo.SetOnline(request.DeviceId)
}

// Heartbeat 设备离线
func (*deviceApp) Heartbeat(_ context.Context, userID, deviceID uint64) error {
	return repo.DeviceRepo.SetOnline(deviceID)
}

// Offline 设备离线
func (*deviceApp) Offline(_ context.Context, deviceID uint64, clientAddr string) error {
	return repo.DeviceRepo.SetOffline(deviceID)
}

// ListByUserID 获取用户所有在线设备
func (*deviceApp) ListByUserID(_ context.Context, userID uint64) ([]domain.Device, error) {
	return repo.DeviceRepo.ListByUserID(userID)
}

// Save 获取设备信息
func (*deviceApp) Save(_ context.Context, pbdevice *pb.Device) (uint64, error) {
	device := &domain.Device{
		ID:            pbdevice.Id,
		Type:          pbdevice.Type,
		Brand:         pbdevice.Brand,
		Model:         pbdevice.Model,
		SystemVersion: pbdevice.SystemVersion,
		SDKVersion:    pbdevice.SdkVersion,
		BrandPushID:   pbdevice.BranchPushId,
	}

	err := repo.DeviceRepo.Save(device)
	return device.ID, err
}

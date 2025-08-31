package device

import (
	"context"

	"gim/pkg/protocol/pb/businesspb"
	pb "gim/pkg/protocol/pb/logicpb"
	"gim/pkg/rpc"
)

type app struct{}

var App = new(app)

// SignIn 登录
func (*app) SignIn(ctx context.Context, request *pb.SignInRequest) error {
	_, err := rpc.GetUserIntClient().Auth(ctx, &businesspb.AuthRequest{
		UserId:   request.UserId,
		DeviceId: request.DeviceId,
		Token:    request.Token,
	})
	if err != nil {
		return err
	}

	device, err := Repo.Get(request.DeviceId)
	if err != nil {
		return err
	}
	device.UserID = request.UserId
	device.ConnectAddr = request.ConnectAddr
	device.ClientAddr = request.ClientAddr
	err = Repo.Save(device)
	if err != nil {
		return err
	}

	return Repo.SetOnline(request.DeviceId)
}

// Heartbeat 设备离线
func (*app) Heartbeat(_ context.Context, userID, deviceID uint64) error {
	return Repo.SetOnline(deviceID)
}

// Offline 设备离线
func (*app) Offline(_ context.Context, deviceID uint64, clientAddr string) error {
	return Repo.SetOffline(deviceID)
}

// ListByUserID 获取用户所有在线设备
func (*app) ListByUserID(_ context.Context, userID uint64) ([]Device, error) {
	return Repo.ListByUserID(userID)
}

// Save 获取设备信息
func (*app) Save(_ context.Context, pbdevice *pb.Device) (uint64, error) {
	device := &Device{
		ID:            pbdevice.Id,
		Type:          pbdevice.Type,
		Brand:         pbdevice.Brand,
		Model:         pbdevice.Model,
		SystemVersion: pbdevice.SystemVersion,
		SDKVersion:    pbdevice.SdkVersion,
		BrandPushID:   pbdevice.BranchPushId,
	}

	err := Repo.Save(device)
	return device.ID, err
}

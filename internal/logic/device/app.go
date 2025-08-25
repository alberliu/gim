package device

import (
	"context"
	"log/slog"
	"time"

	"gim/pkg/gerrors"
	pb "gim/pkg/protocol/pb/logicpb"
	"gim/pkg/protocol/pb/userpb"
	"gim/pkg/rpc"
)

type app struct{}

var App = new(app)

// Register 注册设备
func (*app) Register(ctx context.Context, in *pb.RegisterDeviceRequest) (uint64, error) {
	device := Device{
		Type:          in.Type,
		Brand:         in.Brand,
		Model:         in.Model,
		SystemVersion: in.SystemVersion,
		SDKVersion:    in.SdkVersion,
	}

	// 判断设备信息是否合法
	if !device.IsLegal() {
		return 0, gerrors.ErrBadRequest
	}

	err := Repo.Save(&device)
	if err != nil {
		return 0, err
	}

	return device.ID, nil
}

// SignIn 登录
func (*app) SignIn(ctx context.Context, request *pb.ConnSignInRequest) error {
	_, err := rpc.GetUserIntClient().Auth(ctx, &userpb.AuthRequest{
		UserId:   request.UserId,
		DeviceId: request.DeviceId,
		Token:    request.Token,
	})
	if err != nil {
		return err
	}

	// 标记用户在设备上登录
	device, err := Repo.Get(request.DeviceId)
	if err != nil {
		return err
	}
	device.Online(request.UserId, request.ConnAddr, request.ClientAddr)
	return Repo.Save(device)
}

// Offline 设备离线
func (*app) Offline(ctx context.Context, deviceID uint64, clientAddr string) error {
	device, err := Repo.Get(deviceID)
	if err != nil {
		return err
	}

	if device.ClientAddr != clientAddr {
		return nil
	}
	device.Status = OffLine

	return Repo.Save(device)

}

// ListOnlineByUserID 获取用户所有在线设备
func (*app) ListOnlineByUserID(ctx context.Context, userIDs []uint64) ([]*pb.Device, error) {
	devices, err := Repo.ListOnlineByUserID(userIDs)
	if err != nil {
		return nil, err
	}
	pbDevices := make([]*pb.Device, len(devices))
	for i := range devices {
		pbDevices[i] = devices[i].ToProto()
	}
	return pbDevices, nil
}

// GetDevice 获取设备信息
func (*app) GetDevice(ctx context.Context, deviceId uint64) (*pb.Device, error) {
	device, err := Repo.Get(deviceId)
	if err != nil {
		return nil, err
	}

	return device.ToProto(), nil
}

// ServerStop connect服务停止
func (*app) ServerStop(ctx context.Context, connAddr string) error {
	devices, err := Repo.ListOnlineByConnAddr(connAddr)
	if err != nil {
		return err
	}

	for i := range devices {
		// 因为是异步修改设备转台，要避免设备重连，导致状态不一致
		err = Repo.UpdateStatusOffline(devices[i])
		if err != nil {
			slog.Error("DeviceRepo.Save error", "device", devices[i], "error", err)
		}
		time.Sleep(2 * time.Millisecond)
	}
	return nil
}

package app

import (
	"context"
	devicedomain "gim/internal/logic/domain/device"
	"gim/pkg/gerrors"
	"gim/pkg/pb"
)

type deviceApp struct{}

var DeviceApp = new(deviceApp)

// Register 注册设备
func (*deviceApp) Register(ctx context.Context, in *pb.RegisterDeviceReq) (int64, error) {
	device := devicedomain.Device{
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

	err := devicedomain.DeviceRepo.Save(&device)
	if err != nil {
		return 0, err
	}

	return device.Id, nil
}

// SignIn 登录
func (*deviceApp) SignIn(ctx context.Context, userId, deviceId int64, token string, connAddr string, clientAddr string) error {
	return devicedomain.DeviceService.SignIn(ctx, userId, deviceId, token, connAddr, clientAddr)
}

// Offline 设备离线
func (*deviceApp) Offline(ctx context.Context, deviceId int64, clientAddr string) error {
	device, err := devicedomain.DeviceRepo.Get(deviceId)
	if err != nil {
		return err
	}
	if device == nil {
		return nil
	}

	if device.ClientAddr != clientAddr {
		return nil
	}
	device.Status = devicedomain.DeviceOffLine

	err = devicedomain.DeviceRepo.Save(device)
	if err != nil {
		return err
	}
	return nil
}

func (*deviceApp) ListOnlineByUserId(ctx context.Context, userId int64) ([]*pb.Device, error) {
	return devicedomain.DeviceService.ListOnlineByUserId(ctx, userId)
}

func (*deviceApp) GetDevice(ctx context.Context, deviceId int64) (*pb.Device, error) {
	device, err := devicedomain.DeviceRepo.Get(deviceId)
	return device.ToProto(), err
}

func (*deviceApp) ServerStop(ctx context.Context, connAddr string) error {
	return devicedomain.DeviceService.ServerStop(ctx, connAddr)
}

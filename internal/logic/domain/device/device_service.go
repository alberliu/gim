package device

import (
	"context"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc"
	"time"

	"go.uber.org/zap"
)

type deviceService struct{}

var DeviceService = new(deviceService)

// Register 注册设备
func (*deviceService) Register(ctx context.Context, device *Device) error {
	err := DeviceDao.Save(device)
	if err != nil {
		return err
	}

	return nil
}

// SignIn 长连接登录
func (*deviceService) SignIn(ctx context.Context, userId, deviceId int64, token string, connAddr string, clientAddr string) error {
	_, err := rpc.BusinessIntClient.Auth(ctx, &pb.AuthReq{UserId: userId, DeviceId: deviceId, Token: token})
	if err != nil {
		return err
	}

	// 标记用户在设备上登录
	device, err := DeviceRepo.Get(deviceId)
	if err != nil {
		return err
	}
	if device == nil {
		return nil
	}

	device.Online(userId, connAddr, clientAddr)

	err = DeviceRepo.Save(device)
	if err != nil {
		return err
	}
	return nil
}

// Auth 权限验证
func (*deviceService) Auth(ctx context.Context, userId, deviceId int64, token string) error {
	_, err := rpc.BusinessIntClient.Auth(ctx, &pb.AuthReq{UserId: userId, DeviceId: deviceId, Token: token})
	if err != nil {
		return err
	}
	return nil
}

func (*deviceService) ListOnlineByUserId(ctx context.Context, userId int64) ([]*pb.Device, error) {
	devices, err := DeviceRepo.ListOnlineByUserId(userId)
	if err != nil {
		return nil, err
	}
	pbDevices := make([]*pb.Device, len(devices))
	for i := range devices {
		pbDevices[i] = devices[i].ToProto()
	}
	return pbDevices, nil
}

// ServerStop connect服务停止，需要将连接在这台connect上的设备标记为下线
func (*deviceService) ServerStop(ctx context.Context, connAddr string) error {
	devices, err := DeviceRepo.ListOnlineByConnAddr(connAddr)
	if err != nil {
		return err
	}

	for i := range devices {
		// 因为是异步修改设备转台，要避免设备重连，导致状态不一致
		err = DeviceRepo.UpdateStatusOffline(devices[i])
		if err != nil {
			logger.Logger.Error("DeviceRepo.Save error", zap.Any("device", devices[i]), zap.Error(err))
		}
		time.Sleep(2 * time.Millisecond)
	}
	return nil
}

package device

import (
	"context"
	"time"

	"go.uber.org/zap"

	"gim/pkg/logger"
	"gim/pkg/protocol/pb"
	"gim/pkg/rpc"
)

type service struct{}

var Service = new(service)

// Register 注册设备
func (*service) Register(ctx context.Context, device *Device) error {
	err := Repo.Save(device)
	if err != nil {
		return err
	}

	return nil
}

// SignIn 长连接登录
func (*service) SignIn(ctx context.Context, userId, deviceId int64, token string, connAddr string, clientAddr string) error {
	_, err := rpc.GetBusinessIntClient().Auth(ctx, &pb.AuthReq{UserId: userId, DeviceId: deviceId, Token: token})
	if err != nil {
		return err
	}

	// 标记用户在设备上登录
	device, err := Repo.Get(deviceId)
	if err != nil {
		return err
	}
	if device == nil {
		return nil
	}

	device.Online(userId, connAddr, clientAddr)

	err = Repo.Save(device)
	if err != nil {
		return err
	}
	return nil
}

// Auth 权限验证
func (*service) Auth(ctx context.Context, userId, deviceId int64, token string) error {
	_, err := rpc.GetBusinessIntClient().Auth(ctx, &pb.AuthReq{UserId: userId, DeviceId: deviceId, Token: token})
	if err != nil {
		return err
	}
	return nil
}

func (*service) ListOnlineByUserId(ctx context.Context, userId int64) ([]*pb.Device, error) {
	devices, err := Repo.ListOnlineByUserId(userId)
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
func (*service) ServerStop(ctx context.Context, connAddr string) error {
	devices, err := Repo.ListOnlineByConnAddr(connAddr)
	if err != nil {
		return err
	}

	for i := range devices {
		// 因为是异步修改设备转台，要避免设备重连，导致状态不一致
		err = Repo.UpdateStatusOffline(devices[i])
		if err != nil {
			logger.Logger.Error("DeviceRepo.Save error", zap.Any("device", devices[i]), zap.Error(err))
		}
		time.Sleep(2 * time.Millisecond)
	}
	return nil
}

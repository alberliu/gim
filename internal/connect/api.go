package connect

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "gim/pkg/protocol/pb/connectpb"
)

type ConnIntService struct {
	pb.UnsafeConnectIntServiceServer
}

// PushToDevices 投递消息
func (s *ConnIntService) PushToDevices(ctx context.Context, request *pb.PushToDevicesRequest) (*emptypb.Empty, error) {
	reply := &emptypb.Empty{}

	for _, dm := range request.DeviceMessages {
		conn := GetConn(dm.DeviceId)
		if conn == nil {
			slog.Warn("PushToDevices warn conn not found", "device_id", dm.DeviceId)
			continue
		}

		if conn.DeviceID != dm.DeviceId {
			slog.Warn("PushToDevices warn deviceID not equal", "device_id", dm.DeviceId)
			continue
		}
		conn.SendMessage(dm.Message)
	}
	return reply, nil
}

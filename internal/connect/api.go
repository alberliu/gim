package connect

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "gim/pkg/protocol/pb/connectpb"
)

type ConnectIntService struct {
	pb.UnsafeConnectIntServiceServer
}

// PushToDevices 投递消息
func (*ConnectIntService) PushToDevices(ctx context.Context, request *pb.PushToDevicesRequest) (*emptypb.Empty, error) {
	reply := &emptypb.Empty{}

	for _, dm := range request.DeviceMessages {
		conn := GetConn(dm.DeviceId)
		if conn == nil {
			slog.Warn("push to devices conn not found", "device_id", dm.DeviceId)
			continue
		}

		if conn.DeviceID != dm.DeviceId {
			slog.Warn("push to devices device id mismatch", "device_id", dm.DeviceId)
			continue
		}
		conn.SendMessage(dm.Message)
	}
	return reply, nil
}

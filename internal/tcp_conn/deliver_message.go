package tcp_conn

import (
	"context"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/pb"

	"go.uber.org/zap"
)

func DeliverMessage(ctx context.Context, req *pb.DeliverMessageReq) error {
	// 获取设备对应的TCP连接
	conn, ok := server.GetConn(int(req.Fd))
	if !ok {
		logger.Logger.Warn("GetConn warn", zap.Int64("device_id", req.DeviceId), zap.Int64("df", req.Fd))
		return nil
	}

	data := conn.GetData().(ConnData)
	if data.DeviceId != req.DeviceId {
		logger.Logger.Warn("GetConn warn", zap.Int64("device_id", req.DeviceId), zap.Int64("df", req.Fd))
		return nil
	}

	Handler.Send(conn, pb.PackageType_PT_MESSAGE, grpclib.GetCtxRequstId(ctx), nil, req.Message)
	return nil
}

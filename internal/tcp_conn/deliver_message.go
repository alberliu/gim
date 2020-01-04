package tcp_conn

import (
	"context"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/pb"
)

func DeliverMessage(ctx context.Context, req *pb.DeliverMessageReq) error {
	// 获取设备对应的TCP连接
	conn := load(req.DeviceId)
	if conn == nil {
		logger.Sugar.Warn("ctx id nil")
		return nil
	}

	// 发送消息
	conn.Output(pb.PackageType_PT_MESSAGE, grpclib.GetCtxRequstId(ctx), nil, req.Message)
	return nil
}

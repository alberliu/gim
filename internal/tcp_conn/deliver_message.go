package tcp_conn

import (
	"context"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/pb"

	"go.uber.org/zap"

	"github.com/alberliu/gn"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/status"
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

	send(conn, pb.PackageType_PT_MESSAGE, grpclib.GetCtxRequstId(ctx), nil, req.Message)
	return nil
}

func send(c *gn.Conn, pt pb.PackageType, requestId int64, err error, message proto.Message) {
	var output = pb.Output{
		Type:      pt,
		RequestId: requestId,
	}

	if err != nil {
		status, _ := status.FromError(err)
		output.Code = int32(status.Code())
		output.Message = status.Message()
	}

	if message != nil {
		msgBytes, err := proto.Marshal(message)
		if err != nil {
			logger.Sugar.Error(err)
			return
		}
		output.Data = msgBytes
	}

	outputBytes, err := proto.Marshal(&output)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	err = gn.EncodeToFD(c.GetFd(), outputBytes)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
}

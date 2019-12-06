package server

import (
	"context"
	"gim/logic/service"
	"gim/public/logger"
	"gim/public/pb"
	"gim/public/util"
)

type LogicServerExtServer struct{}

// SendMessage 发送消息
func (*LogicServerExtServer) SendMessage(ctx context.Context, in *pb.SendMessageReq) (*pb.SendMessageResp, error) {
	appId, userId, deviceId, err := util.GetCtxData(ctx)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.SendMessageResp{}, err
	}

	err = service.MessageService.Send(Context(), appId, userId, deviceId, *in)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.SendMessageResp{}, err
	}
	return &pb.SendMessageResp{}, nil
}

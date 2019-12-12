package server

import (
	"context"
	"gim/logic/model"
	"gim/logic/service"
	"gim/public/grpclib"
	"gim/public/logger"
	"gim/public/pb"
)

type LogicServerExtServer struct{}

// SendMessage 发送消息
func (*LogicServerExtServer) SendMessage(ctx context.Context, in *pb.SendMessageReq) (*pb.SendMessageResp, error) {
	appId, err := grpclib.GetCtxAppId(ctx)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	sender := model.Sender{
		AppId:      appId,
		SenderType: pb.SenderType_ST_BUSINESS,
	}
	err = service.MessageService.Send(Context(), sender, *in)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return &pb.SendMessageResp{}, nil
}

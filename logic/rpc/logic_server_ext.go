package rpc

import (
	"context"
	"gim/logic/model"
	"gim/logic/service"
	"gim/public/grpclib"
	"gim/public/pb"
)

type LogicServerExtServer struct{}

// SendMessage 发送消息
func (*LogicServerExtServer) SendMessage(ctx context.Context, in *pb.SendMessageReq) (*pb.SendMessageResp, error) {
	appId, err := grpclib.GetCtxAppId(ctx)
	if err != nil {
		return nil, err
	}

	sender := model.Sender{
		AppId:      appId,
		SenderType: pb.SenderType_ST_BUSINESS,
	}
	err = service.MessageService.Send(ctx, sender, *in)
	if err != nil {
		return nil, err
	}
	return &pb.SendMessageResp{}, nil
}

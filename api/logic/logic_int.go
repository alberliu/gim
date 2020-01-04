package logic

import (
	"context"
	"gim/internal/logic/model"
	"gim/internal/logic/service"
	"gim/pkg/pb"
)

type LogicIntServer struct{}

// SignIn 设备登录
func (*LogicIntServer) SignIn(ctx context.Context, req *pb.SignInReq) (*pb.SignInResp, error) {
	return &pb.SignInResp{}, service.AuthService.SignIn(ctx, req.AppId, req.UserId, req.DeviceId, req.Token, req.ConnAddr)
}

// Sync 设备同步消息
func (*LogicIntServer) Sync(ctx context.Context, req *pb.SyncReq) (*pb.SyncResp, error) {
	messages, err := service.MessageService.ListByUserIdAndSeq(ctx, req.AppId, req.UserId, req.Seq)
	if err != nil {
		return nil, err
	}
	return &pb.SyncResp{Messages: model.MessagesToPB(messages)}, nil
}

// MessageACK 设备收到消息ack
func (*LogicIntServer) MessageACK(ctx context.Context, req *pb.MessageACKReq) (*pb.MessageACKResp, error) {
	return &pb.MessageACKResp{}, service.DeviceAckService.Update(ctx, req.DeviceId, req.DeviceAck)
}

// Offline 设备离线
func (*LogicIntServer) Offline(ctx context.Context, req *pb.OfflineReq) (*pb.OfflineResp, error) {
	return &pb.OfflineResp{}, service.DeviceService.Offline(ctx, req.AppId, req.UserId, req.DeviceId)
}

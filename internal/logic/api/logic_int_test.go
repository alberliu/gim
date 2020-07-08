package api

import (
	"context"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/util"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func getLogicIntClient() pb.LogicIntClient {
	conn, err := grpc.Dial("localhost:50000", grpc.WithInsecure())
	if err != nil {
		logger.Sugar.Error(err)
		return nil
	}
	return pb.NewLogicIntClient(conn)
}

func TestLogicIntServer_SignIn(t *testing.T) {
	token, _ := util.GetToken(1, 1, 1, time.Now().Add(time.Hour).Unix(), util.PublicKey)

	resp, err := getLogicIntClient().SignIn(context.TODO(),
		&pb.SignInReq{
			AppId:    1,
			DeviceId: 1,
			UserId:   1,
			Token:    token,
			ConnAddr: "127.0.0.1:5000",
		})
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	logger.Sugar.Info(resp)
}

func TestLogicIntServer_Sync(t *testing.T) {
	resp, err := getLogicIntClient().Sync(metadata.NewOutgoingContext(context.TODO(), metadata.Pairs("key", "val")),
		&pb.SyncReq{
			AppId:    1,
			UserId:   1,
			DeviceId: 1,
			Seq:      0,
		})
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	logger.Sugar.Info(resp)
}

func TestLogicIntServer_MessageACK(t *testing.T) {
	resp, err := getLogicIntClient().MessageACK(metadata.NewOutgoingContext(context.TODO(), metadata.Pairs("key", "val")),
		&pb.MessageACKReq{
			AppId:       1,
			UserId:      1,
			DeviceId:    1,
			MessageId:   "",
			DeviceAck:   1,
			ReceiveTime: 1,
		})
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	logger.Sugar.Info(resp)
}

func TestLogicIntServer_Offline(t *testing.T) {
	resp, err := getLogicIntClient().Offline(metadata.NewOutgoingContext(context.TODO(), metadata.Pairs("key", "val")),
		&pb.OfflineReq{
			AppId:    1,
			UserId:   1,
			DeviceId: 1,
		})
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	logger.Sugar.Info(resp)
}

package logic

import (
	"context"
	"fmt"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/util"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

func getLogicServerExtClient() pb.LogicServerExtClient {
	conn, err := grpc.Dial("localhost:50002", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return pb.NewLogicServerExtClient(conn)
}

func getServerCtx() context.Context {
	token, _ := util.GetToken(1, 0, 0, time.Now().Add(1*time.Hour).Unix(), util.PublicKey)
	return metadata.NewOutgoingContext(context.TODO(), metadata.Pairs("app_id", "1", "user_id", "0", "device_id", "0", "token", token))
}

func TestLogicServerExtServer_SendMessage(t *testing.T) {
	resp, err := getLogicServerExtClient().SendMessage(getServerCtx(),
		&pb.SendMessageReq{
			MessageId:    "11111",
			ReceiverType: pb.ReceiverType_RT_USER,
			ReceiverId:   1,
			ToUserIds:    nil,
			MessageBody: &pb.MessageBody{
				MessageType: pb.MessageType_MT_TEXT,
				MessageContent: &pb.MessageContent{
					Content: &pb.MessageContent_Text{
						Text: &pb.Text{
							Text: "hello ws",
						},
					},
				},
			},
			IsPersist: true,
			SendTime:  0,
		})
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	logger.Sugar.Info(resp)
}

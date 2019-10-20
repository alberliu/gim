package client

import (
	"gim/conf"
	"gim/public/logger"
	"gim/public/pb"
	"gim/public/transfer"
	"net/rpc"
	"time"

	"github.com/golang/protobuf/proto"
)

var client *rpc.Client

func InitRpcClient() {
	var err error
	for {
		client, err = rpc.Dial("tcp", conf.ConnectRPCServerIP)
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
}

type connectRpcClient struct{}

var ConnectRpcClient = new(connectRpcClient)

// Message 消息投递
func (connectRpcClient) Message(deviceId int64, message pb.Message) (*transfer.MessageResp, error) {
	bytes, err := proto.Marshal(&message)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	req := transfer.MessageReq{
		DeviceId: deviceId,
		Bytes:    bytes,
	}
	var resp = new(transfer.MessageResp)
	err = client.Call("ConnectRPCServer.Message", req, &resp)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return resp, nil
}

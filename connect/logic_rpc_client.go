package connect

import (
	"goim/conf"
	"goim/public/logger"
	"goim/public/transfer"
	"net/rpc"
	"time"
)

var client *rpc.Client

func InitRpcClient() {
	var err error
	for {
		client, err = rpc.Dial("tcp", conf.LogicRPCServerIP)
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
}

type rpcClient struct{}

var RpcClient = new(rpcClient)

// signIn 登录
func (rpcClient) SignIn(req transfer.SignInReq) (*transfer.SignInResp, error) {
	var resp = new(transfer.SignInResp)
	err := client.Call("LogicRPCServer.SignIn", req, &resp)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return resp, nil
}

// sync 同步消息
func (rpcClient) Sync(req transfer.SyncReq) (*transfer.SyncResp, error) {
	var resp = new(transfer.SyncResp)
	err := client.Call("LogicRPCServer.Sync", req, &resp)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return resp, nil
}

// MessageAck 调用逻辑层登录
func (rpcClient) MessageACK(req transfer.MessageAckReq) (*transfer.MessageAckResp, error) {
	var resp = new(transfer.MessageAckResp)
	err := client.Call("LogicRPCServer.MessageACK", req, &resp)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return resp, nil
}

// offline 离线
func (rpcClient) Offline(req transfer.OfflineReq) (*transfer.OfflineResp, error) {
	var resp = new(transfer.OfflineResp)
	err := client.Call("LogicRPCServer.Offline", req, &resp)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return resp, nil
}

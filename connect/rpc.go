package connect

import (
	"goim/conf"
	"goim/public/logger"
	"goim/public/transfer"
	"net/rpc"
)

var clients chan *rpc.Client

func init() {
	clients = make(chan *rpc.Client, len(conf.ConnectRPCClientIPs))
	for _, ip := range conf.ConnectRPCClientIPs {
		client, err := rpc.Dial("tcp", ip)
		if err != nil {
			panic(err)
		}
		clients <- client
	}
}

// signIn 调用逻辑层登录
func signIn(signIn transfer.SignIn) (*transfer.SignInACK, error) {
	client := <-clients
	defer func() { clients <- client }()

	var ack = new(transfer.SignInACK)
	err := client.Call("LogicRPCServer.SignIn", signIn, &ack)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return ack, nil
}

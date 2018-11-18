package main

import (
	"goim/connect"
	"goim/logic/controller"
	"goim/logic/rpc/connect_rpc"
	"goim/logic/rpc/logic_rpc"
)

func init() {
	connect.LogicRPC = logic_rpc.LogicRPC
	connect_rpc.ConnectRPC = connect.ConnectRPC
}

func main() {
	go controller.Engine.Run(":8080")

	conf := connect.Conf{
		Address:      "localhost:50002",
		MaxConnCount: 100,
		AcceptCount:  1,
	}
	server := connect.NewTCPServer(conf)
	server.Start()
}

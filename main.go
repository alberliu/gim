package main

import (
	"fmt"
	"goim/connect"
	"goim/logic/controller"
	"goim/logic/rpc/connect_rpc"
	"goim/logic/rpc/logic_rpc"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	connect.LogicRPC = logic_rpc.LogicRPC
	connect_rpc.ConnectRPC = connect.ConnectRPC
}

func main() {
	go controller.Engine.Run(":8081")
	conf := connect.Conf{
		Address:      "localhost:50002",
		MaxConnCount: 100,
		AcceptCount:  1,
	}
	server := connect.NewTCPServer(conf)
	server.Start()

	websocketServer := connect.WebsocketServer{}
	websocketServer.Start()

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("broker received Signal: ", <-chSig)
}

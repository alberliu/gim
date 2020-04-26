package logic

import (
	"gim/config"
	"gim/pkg/pb"
	"gim/pkg/util"
	"net"

	"google.golang.org/grpc"
)

// StartRpcServer 启动rpc服务
func StartRpcServer() {
	go func() {
		defer util.RecoverPanic()

		intListen, err := net.Listen("tcp", config.LogicConf.RPCIntListenAddr)
		if err != nil {
			panic(err)
		}
		intServer := grpc.NewServer(grpc.UnaryInterceptor(LogicIntInterceptor))
		pb.RegisterLogicIntServer(intServer, &LogicIntServer{})
		err = intServer.Serve(intListen)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer util.RecoverPanic()

		extListen, err := net.Listen("tcp", config.LogicConf.ClientRPCExtListenAddr)
		if err != nil {
			panic(err)
		}
		extServer := grpc.NewServer(grpc.UnaryInterceptor(LogicClientExtInterceptor))
		pb.RegisterLogicClientExtServer(extServer, &LogicClientExtServer{})
		err = extServer.Serve(extListen)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer util.RecoverPanic()

		intListen, err := net.Listen("tcp", config.LogicConf.ServerRPCExtListenAddr)
		if err != nil {
			panic(err)
		}
		intServer := grpc.NewServer(grpc.UnaryInterceptor(LogicServerExtInterceptor))
		pb.RegisterLogicServerExtServer(intServer, &LogicServerExtServer{})
		err = intServer.Serve(intListen)
		if err != nil {
			panic(err)
		}
	}()
}

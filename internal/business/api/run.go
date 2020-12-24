package api

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

		intListen, err := net.Listen("tcp", config.Business.RPCIntListenAddr)
		if err != nil {
			panic(err)
		}
		intServer := grpc.NewServer(grpc.UnaryInterceptor(UserIntInterceptor))
		pb.RegisterBusinessIntServer(intServer, &BusinessIntServer{})
		err = intServer.Serve(intListen)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer util.RecoverPanic()

		extListen, err := net.Listen("tcp", config.Business.RPCExtListenAddr)
		if err != nil {
			panic(err)
		}
		extServer := grpc.NewServer(grpc.UnaryInterceptor(UserExtInterceptor))
		pb.RegisterBusinessExtServer(extServer, &BusinessExtServer{})
		err = extServer.Serve(extListen)
		if err != nil {
			panic(err)
		}
	}()

}

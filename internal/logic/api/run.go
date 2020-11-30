package api

import (
	"gim/config"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/util"
	"net"

	"github.com/alberliu/gn"

	"google.golang.org/grpc"
)

// StartRpcServer 启动rpc服务
func StartRpcServer() {
	go func() {
		defer util.RecoverPanic()

		intListen, err := net.Listen("tcp", config.Logic.RPCIntListenAddr)
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
		gn.SetLogger(logger.Sugar)

		extListen, err := net.Listen("tcp", config.Logic.RPCExtListenAddr)
		if err != nil {
			panic(err)
		}
		extServer := grpc.NewServer(grpc.UnaryInterceptor(LogicExtInterceptor))
		pb.RegisterLogicExtServer(extServer, &LogicExtServer{})
		err = extServer.Serve(extListen)
		if err != nil {
			panic(err)
		}
	}()

}

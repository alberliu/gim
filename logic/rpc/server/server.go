package server

import (
	"gim/conf"
	"log"
	"net"
	"net/rpc"
)

func StartRPCServer() {

	rpc.Register(new(LogicRPCServer))

	tcpAddr, err := net.ResolveTCPAddr("tcp", conf.LogicRPCServerIP)
	if err != nil {
		log.Println(err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		rpc.ServeConn(conn)
	}
}

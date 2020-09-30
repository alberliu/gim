package tcp_conn

import (
	"gim/config"
	"gim/pkg/logger"
	"time"

	"github.com/alberliu/gn"
)

var encoder = gn.NewHeaderLenEncoder(2, 1024)

var server *gn.Server

func StartTCPServer() {
	var err error
	server, err = gn.NewServer(config.TCPConn.TCPListenAddr, &handler{},
		gn.NewHeaderLenDecoder(2, 254),
		gn.WithTimeout(5*time.Minute, 11*time.Minute),
		gn.WithAcceptGNum(10),
		gn.WithIOGNum(100))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	server.Run()
}

package tcp_conn

import (
	"gim/config"
	"gim/pkg/logger"
	"time"

	"github.com/alberliu/gn"
)

var server *gn.Server

func StartTCPServer() {
	var err error
	server, err = gn.NewServer(config.TCPConn.TCPListenAddr, &handler{}, gn.WithReadMaxLen(254), gn.WithTimeout(5*time.Minute, 11*time.Minute))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	server.Run()
}

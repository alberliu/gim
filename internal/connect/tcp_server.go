package connect

import (
	"context"
	"gim/config"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc"
	"time"

	"go.uber.org/zap"

	"github.com/alberliu/gn"
)

var encoder = gn.NewHeaderLenEncoder(2, 1024)

var server *gn.Server

func StartTCPServer() {
	gn.SetLogger(logger.Sugar)

	var err error
	server, err = gn.NewServer(config.Connect.TCPListenAddr, &handler{},
		gn.NewHeaderLenDecoder(2),
		gn.WithReadBufferLen(256),
		gn.WithTimeout(5*time.Minute, 11*time.Minute),
		gn.WithAcceptGNum(10),
		gn.WithIOGNum(100))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	server.Run()
}

type handler struct{}

var Handler = new(handler)

func (*handler) OnConnect(c *gn.Conn) {
	// 初始化连接数据
	conn := &Conn{
		CoonType: CoonTypeTCP,
		TCP:      c,
	}
	c.SetData(conn)
	logger.Logger.Debug("connect:", zap.Int32("fd", c.GetFd()), zap.String("addr", c.GetAddr()))
}

func (h *handler) OnMessage(c *gn.Conn, bytes []byte) {
	conn := c.GetData().(*Conn)
	conn.HandleMessage(bytes)
}

func (*handler) OnClose(c *gn.Conn, err error) {
	conn := c.GetData().(*Conn)
	logger.Logger.Debug("close", zap.String("addr", c.GetAddr()), zap.Int64("user_id", conn.UserId),
		zap.Int64("device_id", conn.DeviceId), zap.Error(err))

	DeleteConn(conn.DeviceId)

	if conn.UserId != 0 {
		_, _ = rpc.LogicIntClient.Offline(context.TODO(), &pb.OfflineReq{
			UserId:     conn.UserId,
			DeviceId:   conn.DeviceId,
			ClientAddr: c.GetAddr(),
		})
	}
}

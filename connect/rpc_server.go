package connect

import (
	"gim/conf"
	"gim/public/logger"
	"gim/public/pb"
	"gim/public/transfer"
	"log"
	"net"
	"net/rpc"
)

type ConnectRPCServer struct{}

// Message 投递消息
func (s *ConnectRPCServer) Message(req transfer.MessageReq, resp *transfer.MessageResp) error {
	// 获取设备对应的TCP连接
	ctx := load(req.DeviceId)
	if ctx == nil {
		logger.Sugar.Error("ctx id nil")
		return nil
	}

	// 发送消息
	err := ctx.Codec.Encode(Package{Code: int(pb.PackageType_PT_MESSAGE), Content: req.Bytes}, WriteDeadline)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

func StartRPCServer() {
	rpc.Register(new(ConnectRPCServer))
	tcpAddr, err := net.ResolveTCPAddr("tcp", conf.ConnectRPCServerIP)
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

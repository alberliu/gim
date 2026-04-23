package connect

import (
	"bufio"
	"log/slog"
	"net"
	"time"

	"gim/pkg/codec"
	"gim/pkg/safe"
)

// StartTCPServer 启动TCP服务器
func StartTCPServer(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	slog.Info("tcp server running")
	go accept(listener)
}

func accept(listener *net.TCPListener) {
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			slog.Error("accept tcp failed", "error", err)
			continue
		}

		err = conn.SetKeepAlive(true)
		if err != nil {
			slog.Error("set keep alive failed", "error", err)
		}

		err = conn.SetNoDelay(true)
		if err != nil {
			slog.Error("set no delay failed", "error", err)
		}

		go handleConn(conn)
	}
}

func handleConn(tcpConn *net.TCPConn) {
	conn := &Conn{
		ConnType: ConnTypeTCP,
		TCP:      tcpConn,
		Reader:   bufio.NewReader(tcpConn),
	}

	var err error
	var buf []byte
	defer safe.RecoverPanic()
	defer func() { conn.Close(err) }()

	for {
		err = conn.TCP.SetReadDeadline(time.Now().Add(ReadDeadline))
		if err != nil {
			return
		}

		buf, err = codec.Decode(conn.Reader)
		if err != nil {
			return
		}

		conn.HandlePacket(buf)
	}
}

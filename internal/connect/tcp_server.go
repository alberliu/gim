package connect

import (
	"bufio"
	"log/slog"
	"net"
	"time"

	"gim/pkg/codec"
	"gim/pkg/util"
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
			slog.Error("acceptTCP error", "error", err)
			continue
		}

		err = conn.SetKeepAlive(true)
		if err != nil {
			slog.Error("setKeepAlive error", "error", err)
		}

		err = conn.SetNoDelay(true)
		if err != nil {
			slog.Error("setNoDelay error", "error", err)
		}

		go handleConn(conn)
	}
}

func handleConn(tcpConn *net.TCPConn) {
	defer util.RecoverPanic()

	conn := &Conn{
		ConnType: ConnTypeTCP,
		TCP:      tcpConn,
		Reader:   bufio.NewReader(tcpConn),
	}

	for {
		err := conn.TCP.SetReadDeadline(time.Now().Add(ReadDeadline))
		if err != nil {
			conn.Close(err)
			return
		}

		buf, err := codec.Decode(conn.Reader)
		if err != nil {
			conn.Close(err)
			return
		}

		conn.HandlePacket(buf)
	}
}

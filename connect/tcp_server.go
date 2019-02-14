package connect

import (
	"fmt"
	"goim/public/logger"
	"net"
	"runtime"
)

// Conf server配置文件
type Conf struct {
	Address      string // 端口
	MaxConnCount int    // 最大连接数
	AcceptCount  int    // 接收建立连接的groutine数量
}

// TCPServer TCP服务器
type TCPServer struct {
	Address      string // 端口
	MaxConnCount int    // 最大连接数
	AcceptCount  int    // 接收建立连接的groutine数量
}

// NewTCPServer 创建TCP服务器
func NewTCPServer(conf Conf) *TCPServer {
	return &TCPServer{
		Address:      conf.Address,
		MaxConnCount: conf.MaxConnCount,
		AcceptCount:  conf.AcceptCount,
	}
}

// Start 启动服务器
func (t *TCPServer) Start() {
	addr, err := net.ResolveTCPAddr("tcp", t.Address)
	if err != nil {
		logger.Sugar.Error(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.Sugar.Error("error listening:", err.Error())
		return
	}
	for i := 0; i < t.AcceptCount; i++ {
		go t.Accept(listener)
	}
	select {}
}

// Accept 接收客户端的TCP长连接的建立
func (t *TCPServer) Accept(listener *net.TCPListener) {
	defer RecoverPanic()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			logger.Sugar.Error(err)
			continue
		}

		err = conn.SetKeepAlive(true)
		if err != nil {
			logger.Sugar.Error(err)
		}

		connContext := NewConnContext(conn)
		go connContext.DoConn()
	}
}

// RecoverPanic 恢复panic
func RecoverPanic() {
	err := recover()
	if err != nil {
		logger.Sugar.Error(GetPanicInfo())
	}

}

// PrintStaStack 打印Panic堆栈信息
func GetPanicInfo() string {
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	return fmt.Sprintf("%s", buf[:n])
}

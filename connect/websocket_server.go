package connect

import (
	"fmt"
	"goim/public/logger"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

type websocketConn interface {
	NextReader() (messageType int, r io.Reader, err error)
	NextWriter(messageType int) (io.WriteCloser, error)
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

type websocketTransport struct {
	sync.Mutex
	socket  websocketConn
	reader  io.Reader
	closing chan bool
}

func (c *websocketTransport) Read(b []byte) (n int, err error) {
	var opCode int
	if c.reader == nil {
		// New message
		var r io.Reader
		for {
			if opCode, r, err = c.socket.NextReader(); err != nil {
				return
			}

			if opCode != websocket.BinaryMessage && opCode != websocket.TextMessage {
				continue
			}

			c.reader = r
			break
		}
	}

	// Read from the reader
	n, err = c.reader.Read(b)
	if err != nil {
		if err == io.EOF {
			c.reader = nil
			err = nil
		}
	}
	return
}

// Write writes data to the connection. It is possible to allow writer to time
// out and return a Error with Timeout() == true after a fixed time limit by
// using SetDeadline and SetWriteDeadline on the websocket.
func (c *websocketTransport) Write(b []byte) (n int, err error) {
	// Serialize write to avoid concurrent write
	c.Lock()
	defer c.Unlock()

	var w io.WriteCloser
	if w, err = c.socket.NextWriter(websocket.BinaryMessage); err == nil {
		if n, err = w.Write(b); err == nil {
			err = w.Close()
		}
	}
	return
}

// Close terminates the connection.
func (c *websocketTransport) Close() error {
	return c.socket.Close()
}

// LocalAddr returns the local network address.
func (c *websocketTransport) LocalAddr() net.Addr {
	return c.socket.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *websocketTransport) RemoteAddr() net.Addr {
	return c.socket.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
func (c *websocketTransport) SetDeadline(t time.Time) (err error) {
	if err = c.socket.SetReadDeadline(t); err == nil {
		err = c.socket.SetWriteDeadline(t)
	}
	return
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
func (c *websocketTransport) SetReadDeadline(t time.Time) error {
	return c.socket.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
func (c *websocketTransport) SetWriteDeadline(t time.Time) error {
	return c.socket.SetWriteDeadline(t)
}

type WebsocketServer struct {
}

func newConn(ws websocketConn) net.Conn {
	conn := &websocketTransport{
		socket:  ws,
		closing: make(chan bool),
	}
	return conn
}

func (*WebsocketServer) Start() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", rootHandler)
	//go http.ListenAndServe("0.0.0.0:8888", nil)
	go http.ListenAndServe(":8888", nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Sugar.Error(err)
	}
	conn := newConn(ws)
	//err = con.SetKeepAlive(true)
	//if err != nil {
	//	logger.Sugar.Error(err)
	//}

	connContext := NewConnContext(conn)
	go connContext.DoConn()
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Println("Could not open file.", err)
	}
	fmt.Fprintf(w, "%s", content)
}

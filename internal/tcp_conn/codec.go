package tcp_conn

import (
	"encoding/binary"
	"errors"
	"gim/pkg/logger"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	ErrIllegalValueLen = errors.New("illegal package length") // 违法的包长度
)

// CodecFactory 编解码器工厂
type CodecFactory struct {
	LenLen             int       // 消息长度字节数组长度
	ReadContentMaxLen  int       // 读缓存区内容最大长度
	WriteContentMaxLen int       // 写缓存区内容最大长度
	ReadBufferPool     sync.Pool // 读缓存内存池
	WriteBufferPool    sync.Pool // 写缓存内存池
}

// NewCodecFactory 创建一个编解码工厂
func NewCodecFactory(lenLen, readContentMaxLen, writeContentMaxLen int) *CodecFactory {
	return &CodecFactory{
		LenLen:             lenLen,
		ReadContentMaxLen:  readContentMaxLen,
		WriteContentMaxLen: writeContentMaxLen,
		ReadBufferPool: sync.Pool{
			New: func() interface{} {
				b := make([]byte, readContentMaxLen+lenLen)
				return b
			},
		},
		WriteBufferPool: sync.Pool{
			New: func() interface{} {
				b := make([]byte, writeContentMaxLen+lenLen)
				return b
			},
		},
	}
}

// Codec 编解码器，用来处理tcp的拆包粘包
type Codec struct {
	f       *CodecFactory
	Conn    net.Conn
	ReadBuf buffer // 读缓冲
}

// GetCodec 创建一个编解码器
func (f *CodecFactory) GetCodec(conn net.Conn) *Codec {
	return &Codec{
		f:       f,
		Conn:    conn,
		ReadBuf: newBuffer(f.ReadBufferPool.Get().([]byte)),
	}
}

// Read 从conn里面读取数据，当conn发生阻塞，这个方法也会阻塞
func (c *Codec) Read() (int, error) {
	return c.ReadBuf.readFromReader(c.Conn)
}

// Decode 解码数据
// Package 代表一个解码包
// bool 标识是否还有可读数据
func (c *Codec) Decode() ([]byte, bool, error) {
	var err error
	// 读取数据长度
	lenBuf, err := c.ReadBuf.seek(0, c.f.LenLen)
	if err != nil {
		return nil, false, nil
	}

	// 读取数据内容
	valueLen := int(binary.BigEndian.Uint16(lenBuf))

	// 数据的字节数组长度大于buffer的长度，返回错误
	if valueLen > c.f.ReadContentMaxLen {
		logger.Logger.Error(ErrIllegalValueLen.Error(), zap.Int("value_len", valueLen))
		return nil, false, ErrIllegalValueLen
	}

	valueBuf, err := c.ReadBuf.read(c.f.LenLen, valueLen)
	if err != nil {
		return nil, false, nil
	}
	return valueBuf, true, nil
}

// Encode 编码数据
func (c *Codec) Encode(bytes []byte, duration time.Duration) error {
	var buffer []byte
	if len(bytes) <= c.f.WriteContentMaxLen {
		bufferCache := c.f.WriteBufferPool.Get().([]byte)
		buffer = bufferCache[0 : c.f.LenLen+len(bytes)]

		defer c.f.WriteBufferPool.Put(bufferCache)
	} else {
		buffer = make([]byte, c.f.LenLen+len(bytes))
	}

	// 将消息长度写入buffer
	binary.BigEndian.PutUint16(buffer[0:c.f.LenLen], uint16(len(bytes)))
	// 将消息内容内容写入buffer
	copy(buffer[c.f.LenLen:], bytes)

	err := c.Conn.SetWriteDeadline(time.Now().Add(duration))
	if err != nil {
		return err
	}
	_, err = c.Conn.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}

// Release 释放编解码器（断开TCP链接，以及归还读缓存区的的内存到内存池）
func (c *Codec) Release() error {
	err := c.Conn.Close()
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	c.f.ReadBufferPool.Put(c.ReadBuf.buf)
	return nil
}

package connect

import (
	"encoding/binary"
	"errors"
	"goim/public/logger"
	"net"
	"time"
)

const (
	TypeLen       = 2                       // 消息类型字节数组长度
	LenLen        = 2                       // 消息长度字节数组长度
	HeadLen       = TypeLen + LenLen        // 消息头部字节数组长度（消息类型字节数组长度+消息长度字节数组长度）
	ContentMaxLen = 1020                    // 消息体最大长度
	BufLen        = ContentMaxLen + HeadLen // 缓冲buffer字节数组长度

)

var (
	ErrIllegalValueLen = errors.New("illegal package length") // 违法的包长度
)

// Codec 编解码器，用来处理tcp的拆包粘包
type Codec struct {
	Conn    net.Conn
	ReadBuf buffer // 读缓冲
}

// newCodec 创建一个编解码器
func NewCodec(conn net.Conn) *Codec {
	return &Codec{
		Conn:    conn,
		ReadBuf: newBuffer(conn, BufLen),
	}
}

// Read 从conn里面读取数据，当conn发生阻塞，这个方法也会阻塞
func (c *Codec) Read() (int, error) {
	return c.ReadBuf.readFromReader()
}

// Decode 解码数据
// Package 代表一个解码包
// bool 标识是否还有可读数据
func (c *Codec) Decode() (*Package, bool, error) {
	var err error
	// 读取数据类型
	typeBuf, err := c.ReadBuf.seek(0, TypeLen)
	if err != nil {
		return nil, false, nil
	}

	// 读取数据长度
	lenBuf, err := c.ReadBuf.seek(TypeLen, HeadLen)
	if err != nil {
		return nil, false, nil
	}

	// 读取数据内容
	valueType := int(binary.BigEndian.Uint16(typeBuf))
	valueLen := int(binary.BigEndian.Uint16(lenBuf))

	// 数据的字节数组长度大于buffer的长度，返回错误
	if valueLen > ContentMaxLen {
		logger.Logger.Error(ErrIllegalValueLen.Error())
		return nil, false, ErrIllegalValueLen
	}

	valueBuf, err := c.ReadBuf.read(HeadLen, valueLen)
	if err != nil {
		return nil, false, nil
	}
	message := Package{Code: valueType, Content: valueBuf}
	return &message, true, nil
}

// Encode 编码数据
func (c *Codec) Encode(pack Package, duration time.Duration) error {
	// 申请写入内存 TODO 可以使用内存池做优化
	buffer := make([]byte, HeadLen+len(pack.Content))

	// 将消息类型写入buffer
	binary.BigEndian.PutUint16(buffer[0:TypeLen], uint16(pack.Code))
	// 将消息长度写入buffer
	binary.BigEndian.PutUint16(buffer[LenLen:HeadLen], uint16(len(pack.Content)))
	// 将消息内容内容写入buffer
	copy(buffer[HeadLen:], pack.Content)

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

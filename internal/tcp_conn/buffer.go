package tcp_conn

import (
	"errors"
	"io"
)

var (
	ErrNotEnough = errors.New("not enough")
)

// buffer 读缓冲区,每个tcp长连接对应一个读缓冲区
type buffer struct {
	buf   []byte // 应用内缓存区
	start int    // 有效字节开始位置
	end   int    // 有效字节结束位置
}

// newBuffer 创建一个缓存区
func newBuffer(bytes []byte) buffer {
	return buffer{bytes, 0, 0}
}

func (b *buffer) len() int {
	return b.end - b.start
}

// grow 将有效的字节前移
func (b *buffer) grow() {
	if b.start == 0 {
		return
	}
	copy(b.buf, b.buf[b.start:b.end])
	b.end -= b.start
	b.start = 0
}

// readFromReader 从reader里面读取数据，如果reader阻塞，会发生阻塞
func (b *buffer) readFromReader(reader io.Reader) (int, error) {
	b.grow()
	n, err := reader.Read(b.buf[b.end:])
	if err != nil {
		return n, err
	}
	b.end += n
	return n, nil
}

// seek 返回n个字节，而不产生移位，如果没有足够字节，返回错误
func (b *buffer) seek(start, end int) ([]byte, error) {
	if b.end-b.start >= end-start {
		buf := b.buf[b.start+start : b.start+end]
		return buf, nil
	}
	return nil, ErrNotEnough
}

// read 舍弃offset个字段，读取n个字段,如果没有足够的字节，返回错误
func (b *buffer) read(offset, limit int) ([]byte, error) {
	if b.len() < offset+limit {
		return nil, ErrNotEnough
	}
	b.start += offset
	buf := b.buf[b.start : b.start+limit]
	b.start += limit
	return buf, nil
}

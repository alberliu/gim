package codec

import (
	"bufio"
	"encoding/binary"
	"errors"
)

func GetUvarintLen(x uint64) int {
	i := 0
	for x >= 0x80 {
		x >>= 7
		i++
	}
	return i + 1
}

func Encode(in []byte) []byte {
	length := GetUvarintLen(uint64(len(in)))
	buf := make([]byte, length+len(in))

	binary.PutUvarint(buf, uint64(len(in)))
	copy(buf[length:], in)
	return buf
}

func Decode(reader *bufio.Reader) ([]byte, error) {
	length, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, length)
	n, err := reader.Read(buf)
	if err != nil {
		return nil, err
	}
	if n != int(length) {
		return nil, errors.New("decode invalid length")
	}
	return buf, nil
}

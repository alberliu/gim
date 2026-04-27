package proxy

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

type BytesCodec struct{}

func (BytesCodec) Marshal(v any) ([]byte, error) {
	vv, ok := v.(*[]byte)
	if ok {
		return *vv, nil
	}

	message, ok := v.(proto.Message)
	if ok {
		return proto.Marshal(message)
	}

	return nil, errors.New("unknown type")
}

func (BytesCodec) Unmarshal(data []byte, v any) error {
	vv, ok := v.(*[]byte)
	if ok {
		*vv = data
		return nil
	}

	message, ok := v.(proto.Message)
	if ok {
		return proto.Unmarshal(data, message)
	}

	return errors.New("unknown type")
}

func (BytesCodec) Name() string {
	return "bytesCodec"
}

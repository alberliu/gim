package util

import (
	"encoding/json"
	"gim/pkg/logger"
	"gim/pkg/pb"

	"github.com/golang/protobuf/proto"
	jsoniter "github.com/json-iterator/go"

	"go.uber.org/zap"
)

func JsonMarshal(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		logger.Logger.Error("json序列化：", zap.Error(err))
	}
	return Bytes2str(bytes)
}

func FormatMessage(messageType pb.MessageType, messageContent []byte) string {
	if messageType == pb.MessageType_MT_UNKNOWN {
		logger.Logger.Error("error message type")
		return "error message type"
	}
	var (
		msg proto.Message
		err error
	)
	switch messageType {
	case pb.MessageType_MT_TEXT:
		msg = &pb.Text{}
		err = proto.Unmarshal(messageContent, msg)
	case pb.MessageType_MT_FACE:
		msg = &pb.Text{}
		err = proto.Unmarshal(messageContent, msg)
	case pb.MessageType_MT_VOICE:
		msg = &pb.Text{}
		err = proto.Unmarshal(messageContent, msg)
	case pb.MessageType_MT_IMAGE:
		msg = &pb.Text{}
		err = proto.Unmarshal(messageContent, msg)
	case pb.MessageType_MT_FILE:
		msg = &pb.Text{}
		err = proto.Unmarshal(messageContent, msg)
	case pb.MessageType_MT_LOCATION:
		msg = &pb.Text{}
		err = proto.Unmarshal(messageContent, msg)
	case pb.MessageType_MT_COMMAND:
		msg = &pb.Text{}
		err = proto.Unmarshal(messageContent, msg)
	case pb.MessageType_MT_CUSTOM:
		msg = &pb.Text{}
		err = proto.Unmarshal(messageContent, msg)
	}

	bytes, err := jsoniter.Marshal(msg)
	if err != nil {
		logger.Sugar.Error(err)
		return ""
	}
	return Bytes2str(bytes)
}

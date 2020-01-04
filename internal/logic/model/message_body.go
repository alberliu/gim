package model

import (
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/util"

	jsoniter "github.com/json-iterator/go"
)

// PBToMessageBody 将pb协议转化为消息体
func PBToMessageBody(pbBody *pb.MessageBody) (int, string) {
	if pbBody.MessageType == pb.MessageType_MT_UNKNOWN {
		logger.Logger.Error("error message type")
		return 0, ""
	}

	var content interface{}
	switch pbBody.MessageType {
	case pb.MessageType_MT_TEXT:
		content = pbBody.MessageContent.GetText()
	case pb.MessageType_MT_FACE:
		content = pbBody.MessageContent.GetFace()
	case pb.MessageType_MT_VOICE:
		content = pbBody.MessageContent.GetVoice()
	case pb.MessageType_MT_IMAGE:
		content = pbBody.MessageContent.GetImage()
	case pb.MessageType_MT_FILE:
		content = pbBody.MessageContent.GetFile()
	case pb.MessageType_MT_LOCATION:
		content = pbBody.MessageContent.GetLocation()
	case pb.MessageType_MT_COMMAND:
		content = pbBody.MessageContent.GetCommand()
	case pb.MessageType_MT_CUSTOM:
		content = pbBody.MessageContent.GetCustom()
	}

	bytes, err := jsoniter.Marshal(content)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, ""
	}

	return int(pbBody.MessageType), util.Bytes2str(bytes)
}

// NewMessageBody 创建一个消息体类型
func NewMessageBody(msgType int, msgContent string) *pb.MessageBody {
	content := pb.MessageContent{}
	switch pb.MessageType(msgType) {
	case pb.MessageType_MT_TEXT:
		var text pb.Text
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &text)
		if err != nil {
			logger.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Text{Text: &text}
	case pb.MessageType_MT_FACE:
		var face pb.Face
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &face)
		if err != nil {
			logger.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Face{Face: &face}
	case pb.MessageType_MT_VOICE:
		var voice pb.Voice
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &voice)
		if err != nil {
			logger.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Voice{Voice: &voice}
	case pb.MessageType_MT_IMAGE:
		var image pb.Image
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &image)
		if err != nil {
			logger.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Image{Image: &image}
	case pb.MessageType_MT_FILE:
		var file pb.File
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &file)
		if err != nil {
			logger.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_File{File: &file}
	case pb.MessageType_MT_LOCATION:
		var location pb.Location
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &location)
		if err != nil {
			logger.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Location{Location: &location}
	case pb.MessageType_MT_COMMAND:
		var command pb.Command
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &command)
		if err != nil {
			logger.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Command{Command: &command}
	case pb.MessageType_MT_CUSTOM:
		var custom pb.Custom
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &custom)
		if err != nil {
			logger.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Custom{Custom: &custom}
	}

	return &pb.MessageBody{
		MessageType:    pb.MessageType(msgType),
		MessageContent: &content,
	}
}

package controller

import (
	"gim/logic/model"
	"gim/logic/service"
	"gim/public/imerror"
	"gim/public/pb"
	"gim/public/util"
	"io/ioutil"
	"net/http"

	"github.com/json-iterator/go"
)

func init() {
	g := Engine.Group("/message")
	g.POST("/send", handler(MessageController{}.Send))
}

type MessageController struct{}

// Send 发送消息
func (MessageController) Send(c *context) {
	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, NewWithError(imerror.WrapErrorWithData(imerror.ErrBadRequest, err)))
		return
	}

	var str interface{}
	jsoniter.Get(bytes, "message_body", "message_content").ToVal(&str)
	body, err := jsoniter.Marshal(str)
	if err != nil {
		c.JSON(http.StatusOK, NewWithError(imerror.WrapErrorWithData(imerror.ErrBadRequest, err)))
		return
	}

	var send model.SendMessage
	err = jsoniter.Unmarshal(bytes, &send)
	if err != nil {
		c.JSON(http.StatusOK, NewWithError(imerror.WrapErrorWithData(imerror.ErrBadRequest, err)))
		return
	}
	send.MessageBody.MessageContent = util.Bytes2str(body)

	pbMessageType := pb.MessageType(send.MessageBody.MessageType)
	if pbMessageType == pb.MessageType_MT_UNKNOWN {
		c.JSON(http.StatusOK, NewWithError(imerror.WrapErrorWithData(imerror.ErrBadRequest, err)))
		return
	}
	send.PbBody = model.NewMessageBody(send.MessageBody.MessageType, send.MessageBody.MessageContent)

	c.response(nil, service.MessageService.Send(Context(), c.appId, c.userId, c.deviceId, send))
}

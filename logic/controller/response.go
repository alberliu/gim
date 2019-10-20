package controller

import (
	"gim/public/imerror"
)

const OK = 200

// 请求返回值
const (
	CodeSuccess = 0 // 成功返回
)

// Response 用户响应数据
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewSuccess(data interface{}) *Response {
	return &Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	}
}

func NewWithError(err *imerror.Error) *Response {
	return &Response{
		Code:    int(err.Code),
		Message: err.Message,
		Data:    err.Data,
	}
}

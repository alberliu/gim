package controller

import (
	"goim/public/imerror"
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

func NewWithBError(err *imerror.LError) *Response {
	return &Response{
		Code:    err.Code,
		Message: err.Message,
		Data:    err.Data,
	}
}

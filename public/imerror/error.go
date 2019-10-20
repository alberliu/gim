package imerror

import (
	"gim/public/pb"
)

const (
	CodeSuccess    = 0         // code成功
	MessageSuccess = "success" // message成功
)

// Error 接入层调用错误
type Error struct {
	Code    pb.ErrCode  // 错误码
	Message string      // 错误信息
	Data    interface{} // 数据
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(code pb.ErrCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func WrapErrorWithData(err *Error, data interface{}) *Error {
	return &Error{
		Code:    err.Code,
		Message: err.Message,
		Data:    data,
	}
}

var (
	ErrUnknown           = NewError(pb.ErrCode_EC_SERVER_UNKNOWN, "error unknown error")              // 服务器未知错误
	ErrUnauthorized      = NewError(pb.ErrCode_EC_UNAUTHORIZED, "error unauthorized")                 // 未登录
	ErrNotInGroup        = NewError(pb.ErrCode_EC_IS_NOT_IN_GROUP, "error not in group")              // 没有在群组内
	ErrDeviceNotBindUser = NewError(pb.ErrCode_EC_DEVICE_NOT_BIND_USER, "error device not bind user") // 没有在群组内
	ErrBadRequest        = NewError(pb.ErrCode_EC_BAD_REQUEST, "error bad request")                   // 请求参数错误
	ErrUserAlreadyExist  = NewError(pb.ErrCode_EC_USER_ALREADY_EXIST, "error user already exist")     // 用户已经存在
	ErrGroupAlreadyExist = NewError(pb.ErrCode_EC_GROUP_ALREADY_EXIST, "error group already exist")   // 群组已经存在
	ErrGroupNotExist     = NewError(pb.ErrCode_EC_GROUP_NOT_EXIST, "error group not exist")           // 群组不存在
)

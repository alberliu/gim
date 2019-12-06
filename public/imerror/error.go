package imerror

import (
	"gim/public/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	ErrUnknown           = status.New(codes.Unknown, "error unknown").Err()                           // 服务器未知错误
	ErrUnauthorized      = newError(pb.ErrCode_EC_UNAUTHORIZED, "error unauthorized")                 // 未登录
	ErrNotInGroup        = newError(pb.ErrCode_EC_IS_NOT_IN_GROUP, "error not in group")              // 没有在群组内
	ErrDeviceNotBindUser = newError(pb.ErrCode_EC_DEVICE_NOT_BIND_USER, "error device not bind user") // 没有在群组内
	ErrBadRequest        = newError(pb.ErrCode_EC_BAD_REQUEST, "error bad request")                   // 请求参数错误
	ErrUserAlreadyExist  = newError(pb.ErrCode_EC_USER_ALREADY_EXIST, "error user already exist")     // 用户已经存在
	ErrGroupAlreadyExist = newError(pb.ErrCode_EC_GROUP_ALREADY_EXIST, "error group already exist")   // 群组已经存在
	ErrGroupNotExist     = newError(pb.ErrCode_EC_GROUP_NOT_EXIST, "error group not exist")           // 群组不存在
	ErrUserNotExist      = newError(pb.ErrCode_EC_USER_NOT_EXIST, "error user not exist")             // 用户不存在
)

func newError(code pb.ErrCode, message string) error {
	return status.New(codes.Code(code), message).Err()
}

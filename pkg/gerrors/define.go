package gerrors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnknown      = status.New(codes.Unknown, "服务器异常").Err() // 服务器未知错误
	ErrUnauthorized = newError(10000, "请重新登录")
	ErrBadRequest   = newError(10001, "请求参数错误")

	ErrBadCode         = newError(10010, "验证码错误")
	ErrNotInGroup      = newError(10011, "用户没有在群组中")
	ErrGroupNotExist   = newError(10013, "群组不存在")
	ErrDeviceNotExist  = newError(10014, "设备不存在")
	ErrAlreadyIsFriend = newError(10015, "对方已经是好友了")
)

func newError(code int, message string) error {
	return status.New(codes.Code(code), message).Err()
}

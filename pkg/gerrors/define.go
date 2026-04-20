package gerrors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnknown      = status.New(codes.Unknown, "服务器异常").Err() // 服务器未知错误
	ErrUnauthorized = newError(10000, "请重新登录")
	ErrBadRequest   = newError(10001, "请求参数错误")
)

func newError(code int, message string) error {
	return status.New(codes.Code(code), message).Err()
}

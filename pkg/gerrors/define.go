package gerrors

import (
	"gim/pkg/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

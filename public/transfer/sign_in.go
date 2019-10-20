package transfer

import (
	"gim/public/imerror"
	"gim/public/logger"
	"gim/public/pb"

	"github.com/golang/protobuf/proto"
)

// SignIn 设备登录
type SignInReq struct {
	ConnectIP string // 连接层RPC监听IP
	Bytes     []byte
}

//  SignInACK 设备登录回执
type SignInResp struct {
	ConnectStatus int    // 连接状态
	AppId         int64  // AppId
	UserId        int64  // 用户id
	DeviceId      int64  // 设备id
	Bytes         []byte // 设备登录响应消息包
}

// NewSignInResp 创建NewSignInResp
func NewSignInResp(code int32, message string, appId, userId, deviceId int64) *SignInResp {
	pbResp := pb.SignInResp{
		Code:    code,
		Message: message,
	}
	bytes, err := proto.Marshal(&pbResp)
	if err != nil {
		logger.Sugar.Error(err)
	}
	connectStatus := ConnectStatusBreak
	if code == imerror.CodeSuccess {
		connectStatus = ConnectStatusOK
	}
	return &SignInResp{
		ConnectStatus: connectStatus,
		AppId:         appId,
		UserId:        userId,
		DeviceId:      deviceId,
		Bytes:         bytes,
	}
}

// ErrorToSignInResp 将error转化成SignInResp
func ErrorToSignInResp(err error, appId, userId, deviceId int64) *SignInResp {
	if err != nil {
		e, ok := err.(*imerror.Error)
		if ok {
			return NewSignInResp(int32(e.Code), e.Message, 0, 0, 0)
		} else {
			return NewSignInResp(int32(imerror.ErrUnknown.Code), imerror.ErrUnknown.Message, 0, 0, 0)
		}
	}
	return NewSignInResp(imerror.CodeSuccess, imerror.MessageSuccess, appId, userId, deviceId)
}

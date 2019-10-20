package transfer

import (
	"goim/public/imerror"
	"goim/public/logger"
	"goim/public/pb"

	"github.com/golang/protobuf/proto"
)

// 同步消息触发
type SyncReq struct {
	AppId    int64  // appId
	UserId   int64  // 用户id
	DeviceId int64  // 设备id
	IsSignIn bool   // 标记用户是否登录成功
	Bytes    []byte // 同步消息字节包
}

// 同步消息触发
type SyncResp struct {
	ConnectStatus int    // 连接状态
	Bytes         []byte // 字节包
}

// NewSyncResp 创建NewSyncResp
func NewSyncResp(code int32, message string, messages []*pb.MessageItem) *SyncResp {
	syncResp := pb.SyncResp{
		Code:     code,
		Message:  "",
		Messages: messages,
	}

	bytes, err := proto.Marshal(&syncResp)
	if err != nil {
		logger.Sugar.Error(err)
	}
	connectStatus := ConnectStatusBreak
	if code == imerror.CodeSuccess {
		connectStatus = ConnectStatusOK
	}
	return &SyncResp{
		ConnectStatus: connectStatus,
		Bytes:         bytes,
	}
}

func ErrorToSyncResp(err error, messages []*pb.MessageItem) *SyncResp {
	if err != nil {
		e, ok := err.(*imerror.Error)
		if ok {
			return NewSyncResp(int32(e.Code), e.Message, nil)
		} else {
			return NewSyncResp(int32(imerror.ErrUnknown.Code), imerror.ErrUnknown.Message, nil)
		}
	}
	return NewSyncResp(imerror.CodeSuccess, imerror.MessageSuccess, messages)
}

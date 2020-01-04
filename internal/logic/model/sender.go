package model

import "gim/pkg/pb"

type Sender struct {
	AppId      int64         // appId
	SenderType pb.SenderType // 发送者类型，1：系统，2：用户，3：业务方
	SenderId   int64         // 发送者id
	DeviceId   int64         // 发送者设备id
}

package domain

import pb "gim/pkg/protocol/pb/logicpb"

type Device struct {
	Type   pb.DeviceType // 设备类型
	Token  string        // token
	Expire int64         // 过期时间
}

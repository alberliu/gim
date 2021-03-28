package model

import "gim/pkg/pb"

type Sender struct {
	SenderType pb.SenderType // 发送者类型，1：系统，2：用户，3：业务方
	SenderId   int64         // 发送者id
	DeviceId   int64         // 发送者设备id
	Nickname   string        // 昵称
	AvatarUrl  string        // 头像
	Extra      string        // 扩展字段
}

func SenderToPB(sender Sender) *pb.Sender {
	return &pb.Sender{
		SenderType: sender.SenderType,
		SenderId:   sender.SenderId,
		AvatarUrl:  sender.AvatarUrl,
		Nickname:   sender.Nickname,
		Extra:      sender.Extra,
	}
}

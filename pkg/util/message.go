package util

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	pb "gim/pkg/protocol/pb/logicpb"
)

var MessagePushes = map[pb.PushCode]proto.Message{
	pb.PushCode_PC_USER_MESSAGE:        &pb.UserMessagePush{},
	pb.PushCode_PC_GROUP_MESSAGE:       &pb.GroupMessagePush{},
	pb.PushCode_PC_ADD_FRIEND:          &pb.AddFriendPush{},
	pb.PushCode_PC_AGREE_ADD_FRIEND:    &pb.AgreeAddFriendPush{},
	pb.PushCode_PC_UPDATE_GROUP:        &pb.UpdateGroupPush{},
	pb.PushCode_PC_ADD_GROUP_MEMBERS:   &pb.AddGroupMembersPush{},
	pb.PushCode_PC_REMOVE_GROUP_MEMBER: &pb.RemoveGroupMemberPush{},
}

func MessageToString(msg *pb.Message) string {
	push, ok := MessagePushes[msg.Code]
	if !ok {
		return fmt.Sprintf("%-5d %-5d %s %s", msg.Code, msg.Seq, "unknown", string(msg.Content))
	}

	_ = proto.Unmarshal(msg.Content, push)
	return fmt.Sprintf("%-5d %-5d %s", msg.Code, msg.Seq, push)

	switch msg.Code {
	case pb.PushCode_PC_USER_MESSAGE:
		p := push.(*pb.UserMessagePush)
		return fmt.Sprintf("%-5d %-5d %v %s", msg.Code, msg.Seq, push, string(p.Content))
	case pb.PushCode_PC_GROUP_MESSAGE:
		p := push.(*pb.GroupMessagePush)
		return fmt.Sprintf("%-5d %-5d %v %s", msg.Code, msg.Seq, push, string(p.Content))
	default:
		return fmt.Sprintf("%-5d %-5d %s", msg.Code, msg.Seq, push)
	}
}

package util

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"gim/pkg/protocol/pb/connectpb"
	pb "gim/pkg/protocol/pb/logicpb"
)

var MessagePushes = map[connectpb.Command]proto.Message{
	connectpb.Command_USER_MESSAGE:        &pb.UserMessagePush{},
	connectpb.Command_GROUP_MESSAGE:       &pb.GroupMessagePush{},
	connectpb.Command_ADD_FRIEND:          &pb.AddFriendPush{},
	connectpb.Command_AGREE_ADD_FRIEND:    &pb.AgreeAddFriendPush{},
	connectpb.Command_UPDATE_GROUP:        &pb.UpdateGroupPush{},
	connectpb.Command_ADD_GROUP_MEMBERS:   &pb.AddGroupMembersPush{},
	connectpb.Command_REMOVE_GROUP_MEMBER: &pb.RemoveGroupMemberPush{},
}

func MessageToString(message *connectpb.Message) string {
	push, ok := MessagePushes[message.Command]
	if !ok {
		return fmt.Sprintf("%-5d %-5d %s %s", message.Code, message.Seq, "unknown", string(message.Content))
	}
	_ = proto.Unmarshal(message.Content, push)

	switch message.Command {
	case connectpb.Command_USER_MESSAGE:
		p := push.(*pb.UserMessagePush)
		return fmt.Sprintf("%-5d %-5d %v %s", message.Code, message.Seq, push, string(p.Content))
	case connectpb.Command_GROUP_MESSAGE:
		p := push.(*pb.GroupMessagePush)
		return fmt.Sprintf("%-5d %-5d %v %s", message.Code, message.Seq, push, string(p.Content))
	default:
		return fmt.Sprintf("%-5d %-5d %s", message.Code, message.Seq, push)
	}
}

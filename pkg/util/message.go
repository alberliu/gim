package util

import (
	"encoding/json"
	"fmt"
	"gim/pkg/protocol/pb"

	"google.golang.org/protobuf/proto"
)

var MessagePushes = map[pb.PushCode]proto.Message{
	pb.PushCode_PC_USER_MESSAGE:        &pb.UserMessagePush{},
	pb.PushCode_PC_GROUP_MESSAGE:       &pb.UserMessagePush{},
	pb.PushCode_PC_ADD_FRIEND:          &pb.AddFriendPush{},
	pb.PushCode_PC_AGREE_ADD_FRIEND:    &pb.AgreeAddFriendPush{},
	pb.PushCode_PC_UPDATE_GROUP:        &pb.UpdateGroupPush{},
	pb.PushCode_PC_ADD_GROUP_MEMBERS:   &pb.AddGroupMembersPush{},
	pb.PushCode_PC_REMOVE_GROUP_MEMBER: &pb.RemoveGroupMemberPush{},
}

func MessageToString(msg *pb.Message) string {
	push, ok := MessagePushes[pb.PushCode(msg.Code)]
	if !ok {
		return fmt.Sprintf("%-5d:%s", msg.Code, "unknown")
	}
	proto.Unmarshal(msg.Content, push)

	switch pb.PushCode(msg.Code) {
	case pb.PushCode_PC_USER_MESSAGE:
		msgPush := push.(*pb.UserMessagePush)
		bytes, _ := json.Marshal(push)
		return fmt.Sprintf("%-5d:%s:%s", msg.Code, string(bytes), string(msgPush.Content))
	default:
		proto.Unmarshal(msg.Content, push)
		bytes, _ := json.Marshal(push)
		return fmt.Sprintf("%-5d:%s", msg.Code, string(bytes))
	}
}

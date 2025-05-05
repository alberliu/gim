package domain

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/protobuf/proto"

	"gim/internal/logic/message"
	"gim/pkg/gerrors"
	"gim/pkg/md"
	pb "gim/pkg/protocol/pb/logicpb"
	"gim/pkg/protocol/pb/userpb"
	"gim/pkg/rpc"
	"gim/pkg/util"
)

const (
	UpdateTypeUpdate = 1
	UpdateTypeDelete = 2
)

// Group 群组
type Group struct {
	ID           uint64      // 群组id
	CreatedAt    time.Time   // 创建时间
	UpdatedAt    time.Time   // 更新时间
	Name         string      // 组名
	AvatarUrl    string      // 头像
	Introduction string      // 群简介
	UserNum      int32       // 群组人数
	Extra        string      // 附加字段
	Members      []GroupUser `gorm:"->"` // 群组成员
}

func (g *Group) ToProto() *pb.Group {
	if g == nil {
		return nil
	}

	return &pb.Group{
		GroupId:      g.ID,
		Name:         g.Name,
		AvatarUrl:    g.AvatarUrl,
		Introduction: g.Introduction,
		UserMum:      g.UserNum,
		Extra:        g.Extra,
		CreateTime:   g.CreatedAt.Unix(),
		UpdateTime:   g.UpdatedAt.Unix(),
	}
}

func CreateGroup(userId uint64, in *pb.GroupCreateRequest) *Group {
	now := time.Now()
	group := &Group{
		Name:         in.Name,
		AvatarUrl:    in.AvatarUrl,
		Introduction: in.Introduction,
		Extra:        in.Extra,
		Members:      make([]GroupUser, 0, len(in.MemberIds)+1),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 创建者添加为管理员
	group.Members = append(group.Members, GroupUser{
		GroupID:    group.ID,
		UserID:     userId,
		MemberType: pb.MemberType_GMT_ADMIN,
	})

	// 其让人添加为成员
	for i := range in.MemberIds {
		group.Members = append(group.Members, GroupUser{
			GroupID:    group.ID,
			UserID:     in.MemberIds[i],
			MemberType: pb.MemberType_GMT_MEMBER,
		})
	}
	return group
}

func (g *Group) GetMemberIDs() []uint64 {
	memberIDs := make([]uint64, 0, len(g.Members))
	for i := range g.Members {
		memberIDs = append(memberIDs, g.Members[i].UserID)
	}
	return memberIDs
}

func (g *Group) PushUpdate(ctx context.Context, userId uint64) error {
	userResp, err := rpc.GetUserIntClient().GetUser(ctx, &userpb.GetUserRequest{UserId: userId})
	if err != nil {
		return err
	}
	return g.PushMessage(ctx, pb.PushCode_PC_UPDATE_GROUP, &pb.UpdateGroupPush{
		OptId:        userId,
		OptName:      userResp.User.Nickname,
		Name:         g.Name,
		AvatarUrl:    g.AvatarUrl,
		Introduction: g.Introduction,
		Extra:        g.Extra,
	}, true)
}

// SendMessage 消息发送至群组
func (g *Group) SendMessage(ctx context.Context, fromDeviceID, fromUserID uint64, req *pb.SendGroupMessageRequest) (uint64, error) {
	if !g.IsMember(fromUserID) {
		slog.Error("SendMessage is not member", "fromUserID", fromUserID, "groupId", req.GroupId)
		return 0, gerrors.ErrNotInGroup
	}

	sender, err := rpc.GetSender(fromDeviceID, fromUserID)
	if err != nil {
		return 0, err
	}

	push := pb.GroupMessagePush{
		Sender:  sender,
		GroupId: req.GroupId,
		Content: req.Content,
	}
	bytes, err := proto.Marshal(&push)
	if err != nil {
		return 0, err
	}

	msg := &pb.Message{
		Code:      pb.PushCode_PC_GROUP_MESSAGE,
		Content:   bytes,
		CreatedAt: util.UnixMilliTime(time.Now()),
	}

	userIDs := make([]uint64, 0, len(g.Members))
	for _, member := range g.Members {
		userIDs = append(userIDs, member.UserID)
	}
	return message.App.SendToUsers(md.NewAndCopyRequestId(ctx), userIDs, msg, true)
}

func (g *Group) IsMember(userId uint64) bool {
	for i := range g.Members {
		if g.Members[i].UserID == userId {
			return true
		}
	}
	return false
}

// PushMessage 向群组推送消息
func (g *Group) PushMessage(ctx context.Context, code pb.PushCode, msg proto.Message, isPersist bool) error {
	go func() {
		defer util.RecoverPanic()
		// 将消息发送给群组用户，使用写扩散
		userIDs := g.GetMemberIDs()

		_, err := message.App.PushToUser(md.NewAndCopyRequestId(ctx), userIDs, code, msg, isPersist)
		if err != nil {
			slog.Error("PushMessage", "error", err)
		}
	}()
	return nil
}

// GetMembers 获取群组用户
func (g *Group) GetMembers(ctx context.Context) ([]*pb.GroupMember, error) {
	members := g.Members
	userIds := make(map[uint64]int32, len(members))
	for i := range members {
		userIds[members[i].UserID] = 0
	}
	resp, err := rpc.GetUserIntClient().GetUsers(ctx, &userpb.GetUsersRequest{UserIds: userIds})
	if err != nil {
		return nil, err
	}

	var infos = make([]*pb.GroupMember, len(members))
	for i, member := range g.Members {
		member := pb.GroupMember{
			UserId:     member.UserID,
			MemberType: pb.MemberType(member.MemberType),
			Remarks:    member.Remarks,
			Extra:      member.Extra,
		}

		user, ok := resp.Users[member.UserId]
		if ok {
			member.Nickname = user.Nickname
			member.Sex = user.Sex
			member.AvatarUrl = user.AvatarUrl
			member.UserExtra = user.Extra
		}
		infos[i] = &member
	}

	return infos, nil
}

// AddMembers 给群组添加用户
func (g *Group) AddMembers(userIds []uint64) ([]GroupUser, error) {
	var members []GroupUser

	for i, userId := range userIds {
		if g.IsMember(userId) {
			return nil, gerrors.ErrUserAlreadyInGroup
		}

		members = append(members, GroupUser{
			GroupID:    g.ID,
			UserID:     userIds[i],
			MemberType: pb.MemberType_GMT_MEMBER,
		})
	}

	g.UserNum += int32(len(userIds))
	return members, nil
}

func (g *Group) PushAddMember(ctx context.Context, optUserId uint64, users []GroupUser) error {
	var addIdMap = make(map[uint64]int32, len(users))
	for i := range users {
		addIdMap[users[i].UserID] = 0
	}

	addIdMap[optUserId] = 0
	usersResp, err := rpc.GetUserIntClient().GetUsers(ctx, &userpb.GetUsersRequest{UserIds: addIdMap})
	if err != nil {
		return err
	}
	var members []*pb.GroupMember
	for k := range addIdMap {
		member, ok := usersResp.Users[k]
		if !ok {
			continue
		}

		members = append(members, &pb.GroupMember{
			UserId:    member.UserId,
			Nickname:  member.Nickname,
			Sex:       member.Sex,
			AvatarUrl: member.AvatarUrl,
			UserExtra: member.Extra,
			Remarks:   "",
			Extra:     "",
		})
	}

	optUser := usersResp.Users[optUserId]
	return g.PushMessage(ctx, pb.PushCode_PC_ADD_GROUP_MEMBERS, &pb.AddGroupMembersPush{
		OptId:   optUser.UserId,
		OptName: optUser.Nickname,
		Members: members,
	}, true)
}

func (g *Group) GetMember(ctx context.Context, userId uint64) *GroupUser {
	for i := range g.Members {
		if g.Members[i].UserID == userId {
			return &g.Members[i]
		}
	}
	return nil
}

func (g *Group) PushDeleteMember(ctx context.Context, optId, userId uint64) error {
	userResp, err := rpc.GetUserIntClient().GetUser(ctx, &userpb.GetUserRequest{UserId: optId})
	if err != nil {
		return err
	}
	return g.PushMessage(ctx, pb.PushCode_PC_REMOVE_GROUP_MEMBER, &pb.RemoveGroupMemberPush{
		OptId:         optId,
		OptName:       userResp.User.Nickname,
		DeletedUserId: userId,
	}, true)
}

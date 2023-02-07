package entity

import (
	"context"
	"gim/internal/logic/proxy"
	"gim/pkg/gerrors"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/protocol/pb"
	"gim/pkg/rpc"
	"gim/pkg/util"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	UpdateTypeUpdate = 1
	UpdateTypeDelete = 2
)

// Group 群组
type Group struct {
	Id           int64       // 群组id
	Name         string      // 组名
	AvatarUrl    string      // 头像
	Introduction string      // 群简介
	UserNum      int32       // 群组人数
	Extra        string      // 附加字段
	CreateTime   time.Time   // 创建时间
	UpdateTime   time.Time   // 更新时间
	Members      []GroupUser `gorm:"-"` // 群组成员
}

type GroupUser struct {
	Id         int64     // 自增主键
	GroupId    int64     // 群组id
	UserId     int64     // 用户id
	MemberType int       // 群组类型
	Remarks    string    // 备注
	Extra      string    // 附加属性
	Status     int       // 状态
	CreateTime time.Time // 创建时间
	UpdateTime time.Time // 更新时间
	UpdateType int       `gorm:"-"` // 更新类型
}

func (g *Group) ToProto() *pb.Group {
	if g == nil {
		return nil
	}

	return &pb.Group{
		GroupId:      g.Id,
		Name:         g.Name,
		AvatarUrl:    g.AvatarUrl,
		Introduction: g.Introduction,
		UserMum:      g.UserNum,
		Extra:        g.Extra,
		CreateTime:   g.CreateTime.Unix(),
		UpdateTime:   g.UpdateTime.Unix(),
	}
}

func CreateGroup(userId int64, in *pb.CreateGroupReq) *Group {
	now := time.Now()
	group := &Group{
		Name:         in.Name,
		AvatarUrl:    in.AvatarUrl,
		Introduction: in.Introduction,
		Extra:        in.Extra,
		Members:      make([]GroupUser, 0, len(in.MemberIds)+1),
		CreateTime:   now,
		UpdateTime:   now,
	}

	// 创建者添加为管理员
	group.Members = append(group.Members, GroupUser{
		GroupId:    group.Id,
		UserId:     userId,
		MemberType: int(pb.MemberType_GMT_ADMIN),
		CreateTime: now,
		UpdateTime: now,
		UpdateType: UpdateTypeUpdate,
	})

	// 其让人添加为成员
	for i := range in.MemberIds {
		group.Members = append(group.Members, GroupUser{
			GroupId:    group.Id,
			UserId:     in.MemberIds[i],
			MemberType: int(pb.MemberType_GMT_MEMBER),
			CreateTime: now,
			UpdateTime: now,
			UpdateType: UpdateTypeUpdate,
		})
	}
	return group
}

func (g *Group) Update(ctx context.Context, in *pb.UpdateGroupReq) error {
	g.Name = in.Name
	g.AvatarUrl = in.AvatarUrl
	g.Introduction = in.Introduction
	g.Extra = in.Extra
	return nil
}

func (g *Group) PushUpdate(ctx context.Context, userId int64) error {
	userResp, err := rpc.GetBusinessIntClient().GetUser(ctx, &pb.GetUserReq{UserId: userId})
	if err != nil {
		return err
	}
	err = g.PushMessage(ctx, pb.PushCode_PC_UPDATE_GROUP, &pb.UpdateGroupPush{
		OptId:        userId,
		OptName:      userResp.User.Nickname,
		Name:         g.Name,
		AvatarUrl:    g.AvatarUrl,
		Introduction: g.Introduction,
		Extra:        g.Extra,
	}, true)
	if err != nil {
		return err
	}
	return nil
}

// SendMessage 消息发送至群组
func (g *Group) SendMessage(ctx context.Context, fromDeviceID, fromUserID int64, req *pb.SendMessageReq) (int64, error) {
	if !g.IsMember(fromUserID) {
		logger.Sugar.Error(ctx, fromUserID, req.ReceiverId, "不在群组内")
		return 0, gerrors.ErrNotInGroup
	}

	sender, err := rpc.GetSender(fromDeviceID, fromUserID)
	if err != nil {
		return 0, err
	}

	push := pb.UserMessagePush{
		Sender:     sender,
		ReceiverId: req.ReceiverId,
		Content:    req.Content,
	}
	bytes, err := proto.Marshal(&push)
	if err != nil {
		return 0, err
	}

	msg := &pb.Message{
		Code:     int32(pb.PushCode_PC_GROUP_MESSAGE),
		Content:  bytes,
		SendTime: req.SendTime,
	}

	// 如果发送者是用户，将消息发送给发送者,获取用户seq
	userSeq, err := proxy.MessageProxy.SendToUser(ctx, fromDeviceID, fromUserID, msg, true)
	if err != nil {
		return 0, err
	}

	go func() {
		defer util.RecoverPanic()
		// 将消息发送给群组用户，使用写扩散
		for _, user := range g.Members {
			// 前面已经发送过，这里不需要再发送
			if user.UserId == sender.UserId {
				continue
			}
			_, err := proxy.MessageProxy.SendToUser(grpclib.NewAndCopyRequestId(ctx), fromDeviceID, user.UserId, msg, true)
			if err != nil {
				return
			}
		}
	}()

	return userSeq, nil
}

func (g *Group) IsMember(userId int64) bool {
	for i := range g.Members {
		if g.Members[i].UserId == userId {
			return true
		}
	}
	return false
}

// PushMessage 向群组推送消息
func (g *Group) PushMessage(ctx context.Context, code pb.PushCode, message proto.Message, isPersist bool) error {
	go func() {
		defer util.RecoverPanic()
		// 将消息发送给群组用户，使用写扩散
		for _, user := range g.Members {
			_, err := proxy.PushToUser(grpclib.NewAndCopyRequestId(ctx), user.UserId, code, message, isPersist)
			if err != nil {
				return
			}
		}
	}()
	return nil
}

// GetMembers 获取群组用户
func (g *Group) GetMembers(ctx context.Context) ([]*pb.GroupMember, error) {
	members := g.Members
	userIds := make(map[int64]int32, len(members))
	for i := range members {
		userIds[members[i].UserId] = 0
	}
	resp, err := rpc.GetBusinessIntClient().GetUsers(ctx, &pb.GetUsersReq{UserIds: userIds})
	if err != nil {
		return nil, err
	}

	var infos = make([]*pb.GroupMember, len(members))
	for i := range members {
		member := pb.GroupMember{
			UserId:     members[i].UserId,
			MemberType: pb.MemberType(members[i].MemberType),
			Remarks:    members[i].Remarks,
			Extra:      members[i].Extra,
		}

		user, ok := resp.Users[members[i].UserId]
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
func (g *Group) AddMembers(ctx context.Context, userIds []int64) ([]int64, []int64, error) {
	var existIds []int64
	var addedIds []int64

	now := time.Now()
	for i, userId := range userIds {
		if g.IsMember(userId) {
			existIds = append(existIds, userIds[i])
			continue
		}

		g.Members = append(g.Members, GroupUser{
			GroupId:    g.Id,
			UserId:     userIds[i],
			MemberType: int(pb.MemberType_GMT_MEMBER),
			CreateTime: now,
			UpdateTime: now,
			UpdateType: UpdateTypeUpdate,
		})
		addedIds = append(addedIds, userIds[i])
	}

	g.UserNum += int32(len(addedIds))

	return existIds, addedIds, nil
}

func (g *Group) PushAddMember(ctx context.Context, optUserId int64, addedIds []int64) error {
	var addIdMap = make(map[int64]int32, len(addedIds))
	for i := range addedIds {
		addIdMap[addedIds[i]] = 0
	}

	addIdMap[optUserId] = 0
	usersResp, err := rpc.GetBusinessIntClient().GetUsers(ctx, &pb.GetUsersReq{UserIds: addIdMap})
	if err != nil {
		return err
	}
	var members []*pb.GroupMember
	for k, _ := range addIdMap {
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
	err = g.PushMessage(ctx, pb.PushCode_PC_ADD_GROUP_MEMBERS, &pb.AddGroupMembersPush{
		OptId:   optUser.UserId,
		OptName: optUser.Nickname,
		Members: members,
	}, true)
	if err != nil {
		logger.Sugar.Error(err)
	}
	return nil
}

func (g *Group) GetMember(ctx context.Context, userId int64) *GroupUser {
	for i := range g.Members {
		if g.Members[i].UserId == userId {
			return &g.Members[i]
		}
	}
	return nil
}

// UpdateMember 更新群组成员信息
func (g *Group) UpdateMember(ctx context.Context, in *pb.UpdateGroupMemberReq) error {
	member := g.GetMember(ctx, in.UserId)
	if member == nil {
		return nil
	}

	member.MemberType = int(in.MemberType)
	member.Remarks = in.Remarks
	member.Extra = in.Extra
	member.UpdateTime = time.Now()
	member.UpdateType = UpdateTypeUpdate
	return nil
}

// DeleteMember 删除用户群组
func (g *Group) DeleteMember(ctx context.Context, userId int64) error {
	member := g.GetMember(ctx, userId)
	if member == nil {
		return nil
	}

	member.UpdateType = UpdateTypeDelete
	return nil
}

func (g *Group) PushDeleteMember(ctx context.Context, optId, userId int64) error {
	userResp, err := rpc.GetBusinessIntClient().GetUser(ctx, &pb.GetUserReq{UserId: optId})
	if err != nil {
		return err
	}
	err = g.PushMessage(ctx, pb.PushCode_PC_REMOVE_GROUP_MEMBER, &pb.RemoveGroupMemberPush{
		OptId:         optId,
		OptName:       userResp.User.Nickname,
		DeletedUserId: userId,
	}, true)
	if err != nil {
		return err
	}
	return nil
}

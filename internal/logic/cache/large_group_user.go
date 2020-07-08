package cache

import (
	"gim/internal/logic/db"
	"gim/internal/logic/model"
	"gim/pkg/gerrors"
	"gim/pkg/logger"
	"gim/pkg/util"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

const (
	LargeGroupUserKey = "large_group_user:"
)

// largeGroupUserCache 聊天室场景群组的用户信息
type largeGroupUserCache struct{}

var LargeGroupUserCache = new(largeGroupUserCache)

func (*largeGroupUserCache) Key(appId, groupId int64) string {
	return LargeGroupUserKey + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(groupId, 10)
}

// Members 获取群组成员
func (c *largeGroupUserCache) Members(appId, groupId int64) ([]model.GroupUser, error) {
	userMap, err := db.RedisCli.HGetAll(c.Key(appId, groupId)).Result()
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	users := make([]model.GroupUser, 0, len(userMap))
	for _, v := range userMap {
		var user model.GroupUser
		err = jsoniter.Unmarshal(util.Str2bytes(v), &user)
		if err != nil {
			logger.Sugar.Error(err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

// IsMember 是否是群组成员
func (c *largeGroupUserCache) IsMember(appId, groupId, userId int64) (bool, error) {
	is, err := db.RedisCli.HExists(c.Key(appId, groupId), strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		return false, gerrors.WrapError(err)
	}

	return is, nil
}

// MembersNum 获取群组成员数
func (c *largeGroupUserCache) MembersNum(appId, groupId int64) (int64, error) {
	membersNum, err := db.RedisCli.HLen(c.Key(appId, groupId)).Result()
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return membersNum, nil
}

// Set 添加群组成员
func (c *largeGroupUserCache) Set(appId, groupId, userId int64, label, extra string) error {
	var user = model.GroupUser{
		AppId:   appId,
		GroupId: groupId,
		UserId:  userId,
		Label:   label,
		Extra:   extra,
	}
	bytes, err := jsoniter.Marshal(user)
	if err != nil {
		return gerrors.WrapError(err)
	}
	_, err = db.RedisCli.HSet(c.Key(user.AppId, user.GroupId), strconv.FormatInt(user.UserId, 10), bytes).Result()
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Del 删除群组成员
func (c *largeGroupUserCache) Del(appId, groupId int64, userId int64) error {
	_, err := db.RedisCli.HDel(c.Key(appId, groupId), strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

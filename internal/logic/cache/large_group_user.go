package cache

import (
	"gim/internal/logic/model"
	"gim/pkg/db"
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

// Members 获取群组成员
func (c *largeGroupUserCache) Members(groupId int64) ([]model.GroupUser, error) {
	userMap, err := db.RedisCli.HGetAll(LargeGroupUserKey + strconv.FormatInt(groupId, 10)).Result()
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
func (c *largeGroupUserCache) IsMember(groupId, userId int64) (bool, error) {
	is, err := db.RedisCli.HExists(LargeGroupUserKey+strconv.FormatInt(groupId, 10), strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		return false, gerrors.WrapError(err)
	}

	return is, nil
}

// MembersNum 获取群组成员数
func (c *largeGroupUserCache) MembersNum(groupId int64) (int64, error) {
	membersNum, err := db.RedisCli.HLen(LargeGroupUserKey + strconv.FormatInt(groupId, 10)).Result()
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return membersNum, nil
}

// Set 添加群组成员
func (c *largeGroupUserCache) Set(groupId, userId int64, remarks, extra string) error {
	var user = model.GroupUser{
		GroupId: groupId,
		UserId:  userId,
		Remarks: remarks,
		Extra:   extra,
	}
	bytes, err := jsoniter.Marshal(user)
	if err != nil {
		return gerrors.WrapError(err)
	}
	_, err = db.RedisCli.HSet(LargeGroupUserKey+strconv.FormatInt(groupId, 10), strconv.FormatInt(user.UserId, 10), bytes).Result()
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Del 删除群组成员
func (c *largeGroupUserCache) Del(groupId int64, userId int64) error {
	_, err := db.RedisCli.HDel(LargeGroupUserKey+strconv.FormatInt(groupId, 10), strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

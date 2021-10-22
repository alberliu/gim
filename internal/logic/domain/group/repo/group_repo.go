package repo

import (
	"gim/internal/logic/domain/group/model"
)

type groupRepo struct{}

var GroupRepo = new(groupRepo)

// Get 获取群组信息
func (*groupRepo) Get(groupId int64) (*model.Group, error) {
	group, err := GroupCache.Get(groupId)
	if err != nil {
		return nil, err
	}
	if group != nil {
		return group, nil
	}

	group, err = GroupDao.Get(groupId)
	if err != nil {
		return nil, err
	}
	members, err := GroupUserRepo.ListUser(groupId)
	if err != nil {
		return nil, err
	}
	group.Members = members

	err = GroupCache.Set(group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// Save 获取群组信息
func (*groupRepo) Save(group *model.Group) error {
	groupId := group.Id
	err := GroupDao.Save(group)
	if err != nil {
		return err
	}

	members := group.Members
	for i := range members {
		members[i].GroupId = group.Id
		if members[i].UpdateType == model.UpdateTypeUpdate {
			err = GroupUserRepo.Save(&(members[i]))
			if err != nil {
				return err
			}
		}
		if members[i].UpdateType == model.UpdateTypeDelete {
			err = GroupUserRepo.Delete(group.Id, members[i].UserId)
			if err != nil {
				return err
			}
		}
	}

	if groupId != 0 {
		err = GroupCache.Del(groupId)
		if err != nil {
			return err
		}
	}
	return nil
}

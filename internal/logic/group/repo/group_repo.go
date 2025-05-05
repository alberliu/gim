package repo

import (
	"gim/internal/logic/group/domain"
)

type groupRepo struct{}

var GroupRepo = new(groupRepo)

// Get 获取群组信息
func (*groupRepo) Get(groupId uint64) (*domain.Group, error) {
	group, err := GroupCache.Get(groupId)
	if err == nil {
		return group, nil
	}

	group, err = GroupDao.Get(groupId)
	if err != nil {
		return nil, err
	}

	err = GroupCache.Set(group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// Save 获取群组信息
func (*groupRepo) Save(group *domain.Group) error {
	err := GroupDao.Save(group)
	if err != nil {
		return err
	}

	return GroupCache.Delete(group.ID)
}

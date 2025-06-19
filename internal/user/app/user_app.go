package app

import (
	"context"
	"time"

	"gim/internal/user/repo"
	pb "gim/pkg/protocol/pb/userpb"
)

type userApp struct{}

var UserApp = new(userApp)

func (*userApp) Get(ctx context.Context, userID uint64) (*pb.User, error) {
	user, err := repo.UserRepo.Get(userID)
	return user.ToProto(), err
}

func (*userApp) Update(ctx context.Context, userID uint64, req *pb.UpdateUserRequest) error {
	u, err := repo.UserRepo.Get(userID)
	if err != nil {
		return err
	}
	if u == nil {
		return nil
	}

	u.Nickname = req.Nickname
	u.Sex = req.Sex
	u.AvatarUrl = req.AvatarUrl
	u.Extra = req.Extra
	u.UpdatedAt = time.Now()

	return repo.UserRepo.Save(u)
}

func (*userApp) GetByIDs(ctx context.Context, userIDs []uint64) (map[uint64]*pb.User, error) {
	users, err := repo.UserRepo.GetByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	pbUsers := make(map[uint64]*pb.User, len(users))
	for i := range users {
		pbUsers[users[i].ID] = users[i].ToProto()
	}
	return pbUsers, nil
}

func (*userApp) Search(ctx context.Context, key string) ([]*pb.User, error) {
	users, err := repo.UserRepo.Search(key)
	if err != nil {
		return nil, err
	}

	pbUsers := make([]*pb.User, len(users))
	for i, v := range users {
		pbUsers[i] = v.ToProto()
	}
	return pbUsers, nil
}

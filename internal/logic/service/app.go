package service

import (
	"context"
	"gim/internal/logic/cache"
	"gim/internal/logic/dao"
	"gim/internal/logic/model"
)

type appService struct{}

var AppService = new(appService)

// Get 注册设备
func (*appService) Get(ctx context.Context, appId int64) (*model.App, error) {
	app, err := cache.AppCache.Get(appId)
	if err != nil {
		return app, nil
	}
	if app != nil {
		return app, nil
	}

	app, err = dao.AppDao.Get(appId)
	if err != nil {
		return app, nil
	}

	if app != nil {
		err = cache.AppCache.Set(app)
		if err != nil {
			return app, nil
		}
	}

	return app, nil
}

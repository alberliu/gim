package service

import (
	"goim/logic/cache"
	"goim/logic/dao"
	"goim/logic/model"
	"goim/public/imctx"
	"goim/public/logger"
)

type appService struct{}

var AppService = new(appService)

// Get 注册设备
func (*appService) Get(ctx *imctx.Context, appId int64) (*model.App, error) {
	app, err := cache.AppCache.Get(appId)
	if err != nil {
		logger.Sugar.Error(err)
		return app, nil
	}
	if app != nil {
		return app, nil
	}

	app, err = dao.AppDao.Get(ctx, appId)
	if err != nil {
		logger.Sugar.Error(err)
		return app, nil
	}

	if app != nil {
		err = cache.AppCache.Set(app)
		if err != nil {
			logger.Sugar.Error(err)
			return app, nil
		}
	}

	return app, nil
}

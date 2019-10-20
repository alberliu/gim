package service

import (
	"encoding/base64"
	"goim/logic/cache"
	"goim/public/imctx"
	"goim/public/imerror"
	"goim/public/logger"
	"goim/public/util"
	"strconv"
	"strings"
	"time"
)

type authService struct{}

var AuthService = new(authService)

// SignIn 长连接登录
func (*authService) SignIn(ctx *imctx.Context, appId, userId, deviceId int64, token string, connectIP string) error {
	// 用户验证
	err := AuthService.VerifyToken(ctx, appId, userId, deviceId, token)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	// 标记用户在设备上登录
	err = DeviceService.Online(ctx, appId, deviceId, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	// 记录下设备在线的主机IP和端口
	err = cache.DeviceIPCache.Set(deviceId, connectIP)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// Auth 验证用户是否登录
func (*authService) Auth(ctx *imctx.Context, appId, userId, deviceId int64, token string) error {
	return AuthService.VerifyToken(ctx, appId, userId, deviceId, token)
}

// VerifySecretKey 对用户秘钥进行校验
func (*authService) VerifyToken(ctx *imctx.Context, appid, userId, diviceId int64, token string) error {
	app, err := AppService.Get(ctx, appid)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	if app == nil {
		return imerror.ErrBadRequest
	}

	bytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		logger.Sugar.Error(err)
		return imerror.ErrUnauthorized
	}
	result, err := util.RsaDecrypt(bytes, util.Str2bytes(app.PtivateKey))
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	//appId:user_id:device_id:expire token格式
	strs := strings.Split(util.Bytes2str(result), ":")
	if strs[0] != strconv.FormatInt(appid, 10) {
		return imerror.ErrUnauthorized
	}
	if strs[1] != strconv.FormatInt(userId, 10) {
		return imerror.ErrUnauthorized
	}
	if strs[2] != strconv.FormatInt(diviceId, 10) {
		return imerror.ErrUnauthorized
	}

	expire, err := strconv.ParseInt(strs[3], 10, 64)
	if err != nil {
		logger.Sugar.Error(err)
		return nil
	}

	if expire < time.Now().Unix() {
		return imerror.ErrUnauthorized
	}
	return nil
}

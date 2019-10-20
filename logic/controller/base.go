package controller

import (
	"goim/logic/service"
	"goim/public/imctx"
	"goim/public/imerror"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var Engine = gin.New()

func init() {
	Engine.Use(handler(verify))

}

const (
	keyAppId    = "app_id"
	keyUserId   = "user_id"
	keyDeviceId = "device_id"
	keyToken    = "token"
)

// verify 权限校验
func verify(c *context) {
	appIdStr := c.GetHeader(keyAppId)
	userIdStr := c.GetHeader(keyUserId)
	deviceIdStr := c.GetHeader(keyDeviceId)
	token := c.GetHeader(keyToken)
	path := c.Request.URL.Path
	if path == "/device" {
		return
	}

	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, NewWithError(imerror.ErrUnauthorized))
		c.Abort()
		return
	}
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, NewWithError(imerror.ErrUnauthorized))
		c.Abort()
		return
	}
	deviceId, err := strconv.ParseInt(deviceIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, NewWithError(imerror.ErrUnauthorized))
		c.Abort()
		return
	}
	err = service.AuthService.Auth(Context(), appId, userId, deviceId, token)
	if err != nil {
		c.JSON(http.StatusOK, NewWithError(imerror.ErrUnauthorized))
		c.Abort()
		return
	}
	c.Keys = make(map[string]interface{}, 2)
	c.Keys[keyAppId] = appId
	c.Keys[keyUserId] = userId
	c.Keys[keyDeviceId] = deviceId

	c.Next()
}

func Context() *imctx.Context {
	return imctx.NewContext()
}

package controller

import (
	"goim/logic/db"
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
	keyDeviceId = "device_id"
	keyUserId   = "user_id"
)

// verify 权限校验
func verify(c *context) {
	deviceIdStr := c.GetHeader("device_id")
	token := c.GetHeader("token")
	path := c.Request.URL.Path
	if path == "/device" {
		return
	}

	deviceId, err := strconv.ParseInt(deviceIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, NewWithBError(imerror.LErrUnauthorized))
		c.Abort()
		return
	}
	userId, err := service.AuthService.Auth(Context(), deviceId, token)
	if err != nil {
		c.JSON(http.StatusOK, NewWithBError(imerror.LErrUnauthorized))
		c.Abort()
		return
	}
	c.Keys = make(map[string]interface{}, 2)
	c.Keys[keyDeviceId] = deviceId
	if path != "/user" && path != "/user/signin" {
		if userId == 0 {
			c.JSON(http.StatusOK, NewWithBError(imerror.LErrDeviceNotBindUser))
			c.Abort()
			return
		}
		c.Keys[keyUserId] = userId
	}
	c.Next()
}

func Context() *imctx.Context {
	return imctx.NewContext(db.Factoty.GetSession())
}

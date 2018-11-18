package controller

import (
	"net/http"

	"goim/public/imerror"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(*context)

func handler(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := new(context)
		context.Context = c
		if deviceId, ok := c.Keys[keyDeviceId]; ok {
			context.deviceId = deviceId.(int64)
		}
		if userId, ok := c.Keys[keyUserId]; ok {
			context.userId = userId.(int64)
		}
		handler(context)
	}
}

type context struct {
	*gin.Context
	deviceId int64 // 设备id
	userId   int64 // 用户id
}

// badParam 参数错误
func (c *context) badParam(err error) {
	c.Context.JSON(http.StatusOK, NewWithBError(imerror.WrapLErrorWithData(imerror.LErrBadRequest, err)))
}

func (c *context) response(data interface{}, err error) {
	if err != nil {
		if berr, ok := err.(*imerror.LError); ok {
			c.JSON(http.StatusOK, NewWithBError(berr))
			return
		}
		c.JSON(http.StatusOK, NewWithBError(imerror.LErrUnknow))
		return
	}
	c.Context.JSON(http.StatusOK, NewSuccess(data))
}

func (c *context) bindJson(value interface{}) error {
	err := c.ShouldBindJSON(value)
	if err != nil {
		c.JSON(http.StatusOK, NewWithBError(imerror.WrapLErrorWithData(imerror.LErrBadRequest, err)))
		return err
	}
	return nil
}

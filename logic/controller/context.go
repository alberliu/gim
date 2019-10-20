package controller

import (
	"net/http"

	"goim/public/imerror"

	"github.com/gin-gonic/gin"
)

// HandlerFunc 自定义handler
type HandlerFunc func(*context)

// handler 将自定义handler转化成gin标准handler
func handler(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := new(context)
		context.Context = c
		if appId, ok := c.Keys[keyAppId]; ok {
			context.appId = appId.(int64)
		}
		if userId, ok := c.Keys[keyUserId]; ok {
			context.userId = userId.(int64)
		}
		if deviceId, ok := c.Keys[keyDeviceId]; ok {
			context.deviceId = deviceId.(int64)
		}

		handler(context)
	}
}

// context 自定义gin上下文
type context struct {
	*gin.Context
	appId    int64 // appId
	userId   int64 // 用户id
	deviceId int64 // 设备id
}

// badParam 参数错误
func (c *context) badParam(err error) {
	c.Context.JSON(http.StatusOK, NewWithError(imerror.WrapErrorWithData(imerror.ErrBadRequest, err)))
}

// response 返回响应体
func (c *context) response(data interface{}, err error) {
	if err != nil {
		if imErr, ok := err.(*imerror.Error); ok {
			c.JSON(http.StatusOK, NewWithError(imErr))
			return
		}
		c.JSON(http.StatusOK, NewWithError(imerror.ErrUnknown))
		return
	}
	c.Context.JSON(http.StatusOK, NewSuccess(data))
}

// bindJson 将json绑定到接口体
func (c *context) bindJson(value interface{}) error {
	err := c.ShouldBindJSON(value)
	if err != nil {
		c.JSON(http.StatusOK, NewWithError(imerror.WrapErrorWithData(imerror.ErrBadRequest, err)))
		return err
	}
	return nil
}

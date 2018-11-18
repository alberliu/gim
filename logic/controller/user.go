package controller

import (
	"goim/logic/model"
	"goim/logic/service"
)

func init() {
	g := Engine.Group("/user")
	g.POST("", handler(UserController{}.Regist))
	g.POST("/signin", handler(UserController{}.SignIn))
	g.POST("/info", handler(UserController{}.UserInfo))
}

type UserController struct{}

// Regist 用户注册
func (UserController) Regist(c *context) {
	var regist model.UserRegist
	if c.bindJson(&regist) != nil {
		return
	}
	c.response(service.UserService.Regist(Context(), c.deviceId, regist))
}

// SignIn 用户登录
func (UserController) SignIn(c *context) {
	var data struct {
		Number   string `json:"number"`
		Password string `json:"password"`
	}
	if c.bindJson(&data) != nil {
		return
	}
	c.response(service.UserService.SignIn(Context(), c.deviceId, data.Number, data.Password))
}

// UserInfo 获取用户信息
func (UserController) UserInfo(c *context) {
	c.response(service.UserService.Get(Context(), c.userId))
}

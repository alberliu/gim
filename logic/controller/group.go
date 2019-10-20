package controller

import (
	"goim/logic/model"
	"goim/logic/service"
	"strconv"
)

func init() {
	g := Engine.Group("/group")
	g.GET("/all", handler(GroupController{}.All))
	g.GET("/one/:group_id", handler(GroupController{}.Get))
	g.POST("", handler(GroupController{}.Create))
	g.PUT("", handler(GroupController{}.Update))
	g.POST("/user", handler(GroupController{}.AddUser))
	g.DELETE("/user", handler(GroupController{}.DeleteUser))
	g.PUT("/user", handler(GroupController{}.UpdateUser))
}

type GroupController struct{}

// Get 获取群组信息
func (GroupController) All(c *context) {
	c.response(service.GroupUserService.ListByUserId(Context(), c.appId, c.userId))
}

// Get 获取群组信息
func (GroupController) Get(c *context) {
	gourpIdStr := c.Param("group_id")
	groupId, err := strconv.ParseInt(gourpIdStr, 10, 64)
	if err != nil {
		c.badParam(err)
		return
	}
	c.response(service.GroupService.Get(Context(), c.appId, groupId))
}

// Create 创建群组
func (GroupController) Create(c *context) {
	var group model.Group
	if c.bindJson(&group) != nil {
		return
	}
	group.AppId = c.appId
	c.response(nil, service.GroupService.Create(Context(), group))
}

// Update 更细群组信息
func (GroupController) Update(c *context) {
	var group model.Group
	if c.bindJson(&group) != nil {
		return
	}
	group.AppId = c.appId
	c.response(nil, service.GroupService.Update(Context(), group))
}

// AddUser 给群组添加用户
func (GroupController) AddUser(c *context) {
	var add model.GroupUserUpdate
	if c.bindJson(&add) != nil {
		return
	}
	c.response(nil, service.GroupService.AddUser(Context(), c.appId, add.GroupId, add.UserId, add.Label, add.Extra))
}

// DeleteUser 从群组删除成员
func (GroupController) DeleteUser(c *context) {
	var update model.GroupUserUpdate
	if c.bindJson(&update) != nil {
		return
	}
	c.response(nil, service.GroupService.DeleteUser(Context(), c.appId, update.GroupId, update.UserId))
}

// UpdateLabel 更新用户群组备注
func (GroupController) UpdateUser(c *context) {
	var update model.GroupUserUpdate
	if c.bindJson(&update) != nil {
		return
	}
	err := service.GroupService.UpdateUser(Context(), c.appId, update.GroupId, c.userId, update.Label, update.Extra)
	c.response(nil, err)
}

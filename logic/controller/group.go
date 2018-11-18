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
	g.POST("", handler(GroupController{}.CreateAndAddUser))
	g.POST("/user", handler(GroupController{}.AddUser))
	g.DELETE("/user", handler(GroupController{}.DeleteUser))
	g.PUT("/user/label", handler(GroupController{}.UpdateLabel))
}

type GroupController struct{}

// Get 获取群组信息
func (GroupController) All(c *context) {
	c.response(service.GroupService.ListByUserId(Context(), c.userId))
}

// Get 获取群组信息
func (GroupController) Get(c *context) {
	gourpIdStr := c.Param("group_id")
	groupId, err := strconv.ParseInt(gourpIdStr, 10, 64)
	if err != nil {
		c.badParam(err)
		return
	}
	c.response(service.GroupService.Get(Context(), groupId))
}

// CreateAndAddUser 创建群组并且添加成员
func (GroupController) CreateAndAddUser(c *context) {
	var data = struct {
		Name    string  `json:"name"`     // 群组名称
		UserIds []int64 `json:"user_ids"` // 群组成员
	}{}
	if c.bindJson(&data) != nil {
		return
	}
	groupId, err := service.GroupService.CreateAndAddUser(Context(), data.Name, data.UserIds)
	c.response(map[string]int64{"id": groupId}, err)
}

// AddUser 给群组添加用户
func (GroupController) AddUser(c *context) {
	var update model.GroupUserUpdate
	if c.bindJson(&update) != nil {
		return
	}
	c.response(nil, service.GroupService.AddUser(Context(), update.GroupId, update.UserIds))
}

// DeleteUser 从群组删除成员
func (GroupController) DeleteUser(c *context) {
	var update model.GroupUserUpdate
	if c.bindJson(&update) != nil {
		return
	}
	c.response(nil, service.GroupService.DeleteUser(Context(), update.GroupId, update.UserIds))
}

// UpdateLabel 更新用户群组备注
func (GroupController) UpdateLabel(c *context) {
	var json struct {
		GroupId int64  `json:"group_id"`
		Label   string `json:"label"`
	}
	if c.bindJson(&json) != nil {
		return
	}
	err := service.GroupService.UpdateLabel(Context(), json.GroupId, c.userId, json.Label)
	c.response(nil, err)
}

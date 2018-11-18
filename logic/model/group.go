package model

// Group 群组
type Group struct {
	Id        int64       `json:"id"`    // 群组id
	Name      string      `json:"name"`  // 组名
	GroupUser []GroupUser `json:"users"` // 群组用户
}

type GroupUserUpdate struct {
	GroupId int64   `json:"group_id"` // 群组名称
	UserIds []int64 `json:"user_ids"` // 群组成员
}

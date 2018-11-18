package model

// GroupUser 群组成员
type GroupUser struct {
	UserId int64  `json:"user_id"` // 用户id
	Label  string `json:"label"`   // 用户标签
	Number string `json:"number"`  // 手机号
	Name   string `json:"name"`    // 昵称
	Sex    int    `json:"sex"`     // 性别，1:男；2:女
	Img    string `json:"img"`     // 用户头像
}

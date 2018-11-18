package model

import "time"

// Device 设备
type Device struct {
	Id            int64     `json:"id"`             // 设备id
	UserId        int64     `json:"user_id"`        // 用户id
	Token         string    `json:"token"`          // 设备登录的token
	Type          int       `json:"type"`           // 设备类型,1:Android；2：IOS；3：Windows; 4：MacOS；5：Web
	Brand         string    `json:"brand"`          // 手机厂商
	Model         string    `json:"model"`          // 机型
	SystemVersion string    `json:"system_version"` // 系统版本
	APPVersion    string    `json:"app_version"`    // APP版本
	Status        int       `json:"state"`          // 在线状态，0：不在线；1：在线
	CreateTime    time.Time `json:"create_time"`    // 创建时间
	UpdateTime    time.Time `json:"update_time"`    // 更新时间
}

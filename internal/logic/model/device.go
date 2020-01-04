package model

import "time"

const (
	DeviceOnLine  = 1 // 设备在线
	DeviceOffLine = 0 // 设备离线
)

// Device 设备
type Device struct {
	Id            int64     `json:"id"`             // 设备id
	DeviceId      int64     `json:"device_id"`      // 设备id
	AppId         int64     `json:"app_id"`         // app_id
	UserId        int64     `json:"user_id"`        // 用户id
	Type          int32     `json:"type"`           // 设备类型,1:Android；2：IOS；3：Windows; 4：MacOS；5：Web
	Brand         string    `json:"brand"`          // 手机厂商
	Model         string    `json:"model"`          // 机型
	SystemVersion string    `json:"system_version"` // 系统版本
	SDKVersion    string    `json:"sdk_version"`    // SDK版本
	Status        int32     `json:"state"`          // 在线状态，0：不在线；1：在线
	ConnAddr      string    `json:"conn_addr"`      // 连接层服务层地址
	CreateTime    time.Time `json:"create_time"`    // 创建时间
	UpdateTime    time.Time `json:"update_time"`    // 更新时间
}

type DeviceToken struct {
	UserId int64
	Token  string
}

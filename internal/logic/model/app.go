package model

import "time"

type App struct {
	Id         int64     // AppId
	Name       string    // 名称
	PrivateKey string    // 私钥
	CreateTime time.Time // 创建时间
	UpdateTime time.Time // 更新时间
}

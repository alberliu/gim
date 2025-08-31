package device

import (
	"time"

	pb "gim/pkg/protocol/pb/logicpb"
)

// Device 设备
type Device struct {
	ID            uint64        // 设备id
	CreatedAt     time.Time     // 创建时间
	UpdatedAt     time.Time     // 更新时间
	UserID        uint64        // 用户id
	Type          pb.DeviceType // 设备类型
	Brand         string        // 手机厂商
	Model         string        // 机型
	SystemVersion string        // 系统版本
	SDKVersion    string        // SDK版本
	BrandPushID   string        // 厂商推送ID
	ConnectAddr   string        // 连接层服务层地址
	ClientAddr    string        // 客户端地址

	IsOnline bool `gorm:"-"` // 是否在线
}

func (d *Device) IsLegal() bool {
	if d.Type == 0 || d.Brand == "" || d.Model == "" ||
		d.SystemVersion == "" || d.SDKVersion == "" {
		return false
	}
	return true
}

func (d *Device) Online(userID uint64, connectAddr string, clientAddr string) {
	d.UserID = userID
	d.ConnectAddr = connectAddr
	d.ClientAddr = clientAddr
}

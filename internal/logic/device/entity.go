package device

import (
	"time"

	pb "gim/pkg/protocol/pb/logicpb"
)

const (
	OnLine  = 1 // 设备在线
	OffLine = 0 // 设备离线
)

// Device 设备
type Device struct {
	ID            uint64        // 设备id
	CreatedAt     time.Time     // 创建时间
	UpdatedAt     time.Time     // 更新时间
	UserId        uint64        // 用户id
	Type          pb.DeviceType // 设备类型
	Brand         string        // 手机厂商
	Model         string        // 机型
	SystemVersion string        // 系统版本
	SDKVersion    string        // SDK版本
	Status        int32         // 在线状态，0：离线；1：在线
	ConnAddr      string        // 连接层服务层地址
	ClientAddr    string        // 客户端地址
}

func (d *Device) ToProto() *pb.Device {
	return &pb.Device{
		DeviceId:      d.ID,
		UserId:        d.UserId,
		Type:          d.Type,
		Brand:         d.Brand,
		Model:         d.Model,
		SystemVersion: d.SystemVersion,
		SdkVersion:    d.SDKVersion,
		Status:        d.Status,
		ConnAddr:      d.ConnAddr,
		ClientAddr:    d.ClientAddr,
		CreateTime:    d.CreatedAt.Unix(),
		UpdateTime:    d.UpdatedAt.Unix(),
	}
}

func (d *Device) IsLegal() bool {
	if d.Type == 0 || d.Brand == "" || d.Model == "" ||
		d.SystemVersion == "" || d.SDKVersion == "" {
		return false
	}
	return true
}

func (d *Device) Online(userId uint64, connAddr string, clientAddr string) {
	d.UserId = userId
	d.ConnAddr = connAddr
	d.ClientAddr = clientAddr
	d.Status = OnLine
}

func (d *Device) Offline(userId uint64, connAddr string, clientAddr string) {
	d.UserId = userId
	d.ConnAddr = connAddr
	d.ClientAddr = clientAddr
	d.Status = OnLine
}

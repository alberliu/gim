package controller

import (
	"goim/logic/model"
	"goim/logic/service"
	"goim/public/imerror"
)

func init() {
	g := Engine.Group("/device")
	g.POST("", handler(DeviceController{}.Regist))
}

type DeviceController struct{}

// Regist 设备注册
func (DeviceController) Regist(c *context) {
	var device model.Device
	if c.ShouldBindJSON(&device) != nil {
		return
	}

	if device.Type == 0 || device.Brand == "" || device.Model == "" ||
		device.SystemVersion == "" || device.APPVersion == "" {
		c.response(nil, imerror.LErrBadRequest)
		return
	}

	id, token, err := service.DeviceService.Regist(Context(), device)
	c.response(map[string]interface{}{"id": id, "token": token}, err)
}

package connect

import (
	"fmt"
	"gim/public/logger"
	"gim/public/pb"
	"gim/public/transfer"
	"gim/public/util"
	"io"
	"net"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	ReadDeadline  = 10 * time.Minute
	WriteDeadline = 10 * time.Second
)

// ConnContext 连接上下文
type ConnContext struct {
	Codec    *Codec // 编解码器
	IsSignIn bool   // 标记连接是否登录成功
	AppId    int64  // AppId
	DeviceId int64  // 设备id
	UserId   int64  // 用户id
}

// Package 消息包
type Package struct {
	Code    int    // 消息类型
	Content []byte // 消息体
}

func NewConnContext(conn *net.TCPConn) *ConnContext {
	codec := NewCodec(conn)
	return &ConnContext{Codec: codec}
}

// DoConn 处理TCP连接
func (c *ConnContext) DoConn() {
	defer util.RecoverPanic()

	c.HandleConnect()

	for {
		err := c.Codec.Conn.SetReadDeadline(time.Now().Add(ReadDeadline))
		if err != nil {
			c.HandleReadErr(err)
			return
		}

		_, err = c.Codec.Read()
		if err != nil {
			c.HandleReadErr(err)
			return
		}

		for {
			message, ok, err := c.Codec.Decode()
			// 解码出错，需要中断连接
			if err != nil {
				logger.Logger.Error(err.Error())
				c.Release()
				return
			}
			if ok {
				c.HandlePackage(message)
				continue
			}
			break
		}
	}
}

// HandleConnect 建立连接
func (c *ConnContext) HandleConnect() {
	logger.Logger.Info("tcp connect")
}

// HandlePackage 处理消息包
func (c *ConnContext) HandlePackage(pack *Package) {
	// 对未登录的用户进行拦截
	if pack.Code != int(pb.PackageType_PT_SIGN_IN_REQ) && c.IsSignIn == false {
		// 应该告诉用户没有登录
		c.Release()
		return
	}

	switch pb.PackageType(pack.Code) {
	case pb.PackageType_PT_SIGN_IN_REQ:
		c.HandlePackageSignIn(pack)
	case pb.PackageType_PT_SYNC_REQ:
		c.HandlePackageSync(pack)
	case pb.PackageType_PT_HEARTBEAT_REQ:
		c.HandlePackageHeartbeat()
	case pb.PackageType_PT_MESSAGE_ACK:
		c.HandlePackageMessageACK(pack)
	}
	return
}

// HandlePackageSignIn 处理登录消息包
func (c *ConnContext) HandlePackageSignIn(pack *Package) {
	resp, err := RpcClient.SignIn(transfer.SignInReq{
		Bytes: pack.Content,
		// TODO 发送本机IP
	})
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	if resp.ConnectStatus == transfer.ConnectStatusOK {
		// 将连接保存到本机字典
		c.IsSignIn = true
		c.AppId = resp.AppId
		c.DeviceId = resp.DeviceId
		c.UserId = resp.UserId
		store(c.DeviceId, c)
	}
	// todo 登录失败处理

	// 将响应写回缓冲区
	err = c.Codec.Encode(Package{Code: int(pb.PackageType_PT_SIGN_IN_RESP), Content: resp.Bytes}, WriteDeadline)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
}

// HandlePackageSyncTrigger 处理同步触发消息包
func (c *ConnContext) HandlePackageSync(pack *Package) {
	resp, err := RpcClient.Sync(transfer.SyncReq{
		IsSignIn: c.IsSignIn,
		AppId:    c.AppId,
		DeviceId: c.DeviceId,
		UserId:   c.UserId,
		Bytes:    pack.Content,
	})
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	// 将响应写回缓冲区
	err = c.Codec.Encode(Package{Code: int(pb.PackageType_PT_SYNC_RESP), Content: resp.Bytes}, WriteDeadline)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
}

// HandlePackageHeadbeat 处理心跳包
func (c *ConnContext) HandlePackageHeartbeat() {
	err := c.Codec.Encode(Package{Code: int(pb.PackageType_PT_HEARTBEAT_RESP), Content: []byte{}}, WriteDeadline)
	if err != nil {
		logger.Sugar.Error(err)
	}
	logger.Sugar.Infow("heartbeat", "device_id", c.DeviceId, "user_id", c.UserId)
}

// HandlePackageMessageACK 处理消息回执消息包
func (c *ConnContext) HandlePackageMessageACK(pack *Package) {
	resp, err := RpcClient.MessageACK(transfer.MessageAckReq{
		IsSignIn: c.IsSignIn,
		AppId:    c.AppId,
		DeviceId: c.DeviceId,
		UserId:   c.UserId,
		Bytes:    pack.Content,
	})
	if err != nil {
		logger.Sugar.Error(err)
		// todo 响应未知异常
		return
	}
	fmt.Println(resp)
}

// HandleReadErr 读取conn错误
func (c *ConnContext) HandleReadErr(err error) {
	logger.Logger.Debug("read tcp error：", zap.Int64("app_id", c.AppId), zap.Int64("user_id", c.UserId),
		zap.Int64("device_id", c.DeviceId), zap.Error(err))
	str := err.Error()
	// 服务器主动关闭连接
	if strings.HasSuffix(str, "use of closed network connection") {
		return
	}

	c.Release()
	// 客户端主动关闭连接或者异常程序退出
	if err == io.EOF {
		return
	}
	// SetReadDeadline 之后，超时返回的错误
	if strings.HasSuffix(str, "i/o timeout") {
		return
	}
}

// Release 释放TCP连接
func (c *ConnContext) Release() {
	// 从本地manager中删除tcp连接
	delete(c.DeviceId)

	// 关闭tcp连接
	err := c.Codec.Release()
	if err != nil {
		logger.Sugar.Error(err)
	}

	// 通知业务服务器设备下线
	_, err = RpcClient.Offline(transfer.OfflineReq{
		DeviceId: c.DeviceId,
		UserId:   c.UserId,
	})
	if err != nil {
		logger.Sugar.Error(err)
	}
}

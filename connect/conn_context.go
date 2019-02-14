package connect

import (
	"fmt"
	"goim/public/logger"
	"goim/public/pb"
	"goim/public/transfer"
	"io"
	"net"
	"strings"
	"time"

	"goim/public/lib"

	"goim/conf"

	"github.com/golang/protobuf/proto"
)

const (
	ReadDeadline  = 10 * time.Minute
	WriteDeadline = 10 * time.Second
)

// 消息协议
const (
	CodeSignIn         = 1 // 设备登录
	CodeSignInACK      = 2 // 设备登录回执
	CodeSyncTrigger    = 3 // 消息同步触发
	CodeHeadbeat       = 4 // 心跳
	CodeHeadbeatACK    = 5 // 心跳回执
	CodeMessageSend    = 6 // 消息发送
	CodeMessageSendACK = 7 // 消息发送回执
	CodeMessage        = 8 // 消息投递
	CodeMessageACK     = 9 // 消息投递回执
)

// ConnContext 连接上下文
type ConnContext struct {
	Codec    *Codec // 编解码器
	IsSignIn bool   // 是否登录
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
	defer RecoverPanic()

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
			message, ok := c.Codec.Decode()
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
	// 未登录拦截
	if pack.Code != CodeSignIn && c.IsSignIn == false {
		c.Release()
		return
	}

	switch pack.Code {
	case CodeSignIn:
		c.HandlePackageSignIn(pack)
	case CodeSyncTrigger:
		c.HandlePackageSyncTrigger(pack)
	case CodeHeadbeat:
		c.HandlePackageHeadbeat()
	case CodeMessageSend:
		c.HandlePackageMessageSend(pack)
	case CodeMessageACK:
		c.HandlePackageMessageACK(pack)
	}
	return
}

// HandlePackageSignIn 处理登录消息包
func (c *ConnContext) HandlePackageSignIn(pack *Package) {
	var sign pb.SignIn
	err := proto.Unmarshal(pack.Content, &sign)
	if err != nil {
		logger.Sugar.Error(err)
		c.Release()
		return
	}

	transferSignIn := transfer.SignIn{
		DeviceId: sign.DeviceId,
		UserId:   sign.UserId,
		Token:    sign.Token,
	}

	// 处理设备登录逻辑
	ack, err := signIn(transferSignIn)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	content, err := proto.Marshal(&pb.SignInACK{Code: int32(ack.Code), Message: ack.Message})
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	err = c.Codec.Eecode(Package{Code: CodeSignInACK, Content: content}, WriteDeadline)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	if ack.Code == transfer.CodeSignInSuccess {
		// 将连接保存到本机字典
		c.IsSignIn = true
		c.DeviceId = sign.DeviceId
		c.UserId = sign.UserId
		store(c.DeviceId, c)

		// 将设备和服务器IP的对应关系保存到redis
		redisClient.Set(deviceIdPre+fmt.Sprint(c.DeviceId), conf.ConnectTCPListenIP+"."+conf.ConnectTCPListenPort,
			0)
	}
}

// HandlePackageSyncTrigger 处理同步触发消息包
func (c *ConnContext) HandlePackageSyncTrigger(pack *Package) {
	var trigger pb.SyncTrigger
	err := proto.Unmarshal(pack.Content, &trigger)
	if err != nil {
		logger.Sugar.Error(err)
		c.Release()
		return
	}

	transferTrigger := transfer.SyncTrigger{
		DeviceId:     c.DeviceId,
		UserId:       c.UserId,
		SyncSequence: trigger.SyncSequence,
	}

	publishSyncTrigger(transferTrigger)
}

// HandlePackageHeadbeat 处理心跳包
func (c *ConnContext) HandlePackageHeadbeat() {
	err := c.Codec.Eecode(Package{Code: CodeHeadbeatACK, Content: []byte{}}, WriteDeadline)
	if err != nil {
		logger.Sugar.Error(err)
	}
	logger.Sugar.Infow("心跳：", "device_id", c.DeviceId, "user_id", c.UserId)
}

// HandlePackageMessageSend 处理消息发送包
func (c *ConnContext) HandlePackageMessageSend(pack *Package) {
	var send pb.MessageSend
	err := proto.Unmarshal(pack.Content, &send)
	if err != nil {
		logger.Sugar.Error(err)
		c.Release()
		return
	}

	transferSend := transfer.MessageSend{
		SenderDeviceId: c.DeviceId,
		SenderUserId:   c.UserId,
		ReceiverType:   send.ReceiverType,
		ReceiverId:     send.ReceiverId,
		Type:           send.Type,
		Content:        send.Content,
		SendSequence:   send.SendSequence,
		SendTime:       lib.UnunixTime(send.SendTime),
	}

	publishMessageSend(transferSend)
	if err != nil {
		logger.Sugar.Error(err)
	}
}

// HandlePackageMessageACK 处理消息回执消息包
func (c *ConnContext) HandlePackageMessageACK(pack *Package) {
	var ack pb.MessageACK
	err := proto.Unmarshal(pack.Content, &ack)
	if err != nil {
		logger.Sugar.Error(err)
		c.Release()
		return
	}

	transferAck := transfer.MessageACK{
		MessageId:    ack.MessageId,
		DeviceId:     c.DeviceId,
		UserId:       c.UserId,
		SyncSequence: ack.SyncSequence,
		ReceiveTime:  lib.UnunixTime(ack.ReceiveTime),
	}

	publishMessageACK(transferAck)
}

// HandleReadErr 读取conn错误
func (c *ConnContext) HandleReadErr(err error) {
	logger.Sugar.Infow("连接读取异常：", "device_id", c.DeviceId, "user_id", c.UserId, "err_msg", err)
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
	logger.Sugar.Infow("连接读取未知异常：", "device_id", c.DeviceId, "user_id", c.UserId, "err_msg", err)
}

// Release 释放TCP连接
func (c *ConnContext) Release() {
	delete(c.DeviceId)
	err := c.Codec.Conn.Close()
	if err != nil {
		logger.Sugar.Error(err)
	}

	publishOffLine(transfer.OffLine{
		DeviceId: c.DeviceId,
		UserId:   c.UserId,
	})
}

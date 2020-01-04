package tcp_conn

import (
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/util"
	"io"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"

	"go.uber.org/zap"
)

const (
	ReadDeadline  = 10 * time.Minute
	WriteDeadline = 10 * time.Second
)

const (
	TypeLen            = 2   // 消息类型字节数组长度
	LenLen             = 2   // 消息长度字节数组长度
	ReadContentMaxLen  = 252 // 读缓存区内容最大长度
	WriteContentMaxLen = 508 // 写缓存区内容最大长度
)

var codecFactory = NewCodecFactory(LenLen, ReadContentMaxLen, WriteContentMaxLen)

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
	codec := codecFactory.GetCodec(conn)
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
			bytes, ok, err := c.Codec.Decode()
			// 解码出错，需要中断连接
			if err != nil {
				logger.Logger.Error(err.Error())
				c.Release()
				return
			}
			if ok {
				c.HandlePackage(bytes)
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
func (c *ConnContext) HandlePackage(pack []byte) {
	Handler.Handler(c, pack)
}

// Output
func (c *ConnContext) Output(pt pb.PackageType, requestId int64, err error, message proto.Message) {
	var output = pb.Output{
		Type:      pt,
		RequestId: requestId,
	}

	if err != nil {
		status, _ := status.FromError(err)
		output.Code = int32(status.Code())
		output.Message = status.Message()
	}

	if message != nil {
		msgBytes, err := proto.Marshal(message)
		if err != nil {
			logger.Sugar.Error(err)
			return
		}
		output.Data = msgBytes
	}

	outputBytes, err := proto.Marshal(&output)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	err = c.Codec.Encode(outputBytes, WriteDeadline)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
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
	if c.DeviceId != PreConn {
		delete(c.DeviceId)
	}

	// 关闭tcp连接
	err := c.Codec.Release()
	if err != nil {
		logger.Sugar.Error(err)
	}

	// 通知业务服务器设备下线
	if c.DeviceId != PreConn {
		Handler.Offline(c)
	}
}

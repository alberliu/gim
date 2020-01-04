package main

import (
	"fmt"
	"gim/internal/tcp_conn"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/util"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
)

func main() {
	for i := 0; i < 1000; i++ {
		client := TcpClient{
			AppId:    1,
			UserId:   int64(i),
			DeviceId: int64(i),
			Seq:      0,
			codec:    nil,
		}
		client.Start()
	}
	select {}
}

func Json(i interface{}) string {
	bytes, _ := jsoniter.Marshal(i)
	return string(bytes)
}

var codecFactory = tcp_conn.NewCodecFactory(2, 65536, 1024)

type TcpClient struct {
	AppId    int64
	UserId   int64
	DeviceId int64
	Seq      int64
	codec    *tcp_conn.Codec
}

func (c *TcpClient) Output(pt pb.PackageType, requestId int64, message proto.Message) {
	var input = pb.Input{
		Type:      pt,
		RequestId: requestId,
	}

	if message != nil {
		bytes, err := proto.Marshal(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		input.Data = bytes
	}

	inputByf, err := proto.Marshal(&input)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = c.codec.Encode(inputByf, time.Second)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *TcpClient) Start() {
	connect, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	c.codec = codecFactory.GetCodec(connect)

	c.SignIn()
	c.SyncTrigger()
	go c.Heartbeat()
}

func (c *TcpClient) SignIn() {
	token, err := util.GetToken(c.AppId, c.UserId, c.DeviceId, time.Now().Add(24*30*time.Hour).Unix(), util.PublicKey)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	signIn := pb.SignInInput{
		AppId:    c.AppId,
		UserId:   c.UserId,
		DeviceId: c.DeviceId,
		Token:    token,
	}
	c.Output(pb.PackageType_PT_SIGN_IN, time.Now().UnixNano(), &signIn)
}

func (c *TcpClient) SyncTrigger() {
	c.Output(pb.PackageType_PT_SYNC, time.Now().UnixNano(), &pb.SyncInput{Seq: c.Seq})
}

func (c *TcpClient) Heartbeat() {
	ticker := time.NewTicker(time.Minute * 5)
	for range ticker.C {
		c.Output(pb.PackageType_PT_HEARTBEAT, time.Now().UnixNano(), nil)
	}
}

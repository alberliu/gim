package main

import (
	"fmt"
	"gim/pkg/pb"
	"gim/pkg/util"
	"net"
	"time"

	util2 "github.com/alberliu/gn/test/util"
	"github.com/golang/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
)

func main() {
	client := TcpClient{}
	fmt.Println("input UserId,DeviceId,SyncSequence")
	fmt.Scanf("%d %d %d", &client.UserId, &client.DeviceId, &client.Seq)
	client.Start()
	select {}
}

func Json(i interface{}) string {
	bytes, _ := jsoniter.Marshal(i)
	return string(bytes)
}

type TcpClient struct {
	UserId   int64
	DeviceId int64
	Seq      int64
	codec    *util2.Codec
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

	_, err = c.codec.Conn.Write(util2.Encode(inputByf))
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

	c.codec = util2.NewCodec(connect)

	c.SignIn()
	c.SyncTrigger()
	go c.Heartbeat()
	go c.Receive()
}

func (c *TcpClient) SignIn() {
	signIn := pb.SignInInput{
		UserId:   c.UserId,
		DeviceId: c.DeviceId,
		Token:    "0",
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

func (c *TcpClient) Receive() {
	for {
		_, err := c.codec.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			bytes, ok, err := c.codec.Decode()
			if err != nil {
				fmt.Println(err)
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

func (c *TcpClient) HandlePackage(bytes []byte) {
	var output pb.Output
	err := proto.Unmarshal(bytes, &output)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch output.Type {
	case pb.PackageType_PT_SIGN_IN:
		fmt.Println(Json(output))
	case pb.PackageType_PT_HEARTBEAT:
		fmt.Println("心跳响应")
	case pb.PackageType_PT_SYNC:
		fmt.Println("离线消息同步开始------")
		syncResp := pb.SyncOutput{}
		err := proto.Unmarshal(output.Data, &syncResp)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("离线消息同步响应:code", output.Code, "message:", output.Message)
		for _, msg := range syncResp.Messages {
			fmt.Printf("消息：发送者类型：%d 发送者id：%d 请求id：%d 接收者类型：%d 接收者id：%d  消息内容：%+v seq：%d \n",
				msg.SenderType, msg.SenderId, msg.RequestId, msg.ReceiverType, msg.ReceiverId, util.FormatMessage(msg.MessageType, msg.MessageContent), msg.Seq)
			c.Seq = msg.Seq
		}

		ack := pb.MessageACK{
			DeviceAck:   c.Seq,
			ReceiveTime: util.UnixMilliTime(time.Now()),
		}
		c.Output(pb.PackageType_PT_MESSAGE, output.RequestId, &ack)
		fmt.Println("离线消息同步结束------")
	case pb.PackageType_PT_MESSAGE:
		message := pb.Message{}
		err := proto.Unmarshal(output.Data, &message)
		if err != nil {
			fmt.Println(err)
			return
		}

		msg := message.Message
		fmt.Printf("消息：发送者类型：%d 发送者id：%d 请求id：%d 接收者类型：%d 接收者id：%d  消息内容：%+v seq：%d \n",
			msg.SenderType, msg.SenderId, msg.RequestId, msg.ReceiverType, msg.ReceiverId, util.FormatMessage(msg.MessageType, msg.MessageContent), msg.Seq)

		c.Seq = msg.Seq
		ack := pb.MessageACK{
			DeviceAck:   msg.Seq,
			ReceiveTime: util.UnixMilliTime(time.Now()),
		}
		c.Output(pb.PackageType_PT_MESSAGE, output.RequestId, &ack)
	default:
		fmt.Println("switch other")
	}
}

package main

import (
	"fmt"
	"gim/pkg/protocol/pb"
	"gim/pkg/util"
	"log"
	"net"
	"time"

	gim_util "github.com/alberliu/gn/util"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/proto"
)

var codecFactory = gim_util.NewHeaderLenCodecFactory(2, 65536)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	client := TcpClient{}
	log.Println("input UserId,DeviceId,SyncSeq")
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
	codec    *gim_util.Codec
	Conn     net.Conn
}

func (c *TcpClient) Output(pt pb.PackageType, requestId int64, message proto.Message) {
	var input = pb.Input{
		Type:      pt,
		RequestId: requestId,
	}

	if message != nil {
		bytes, err := proto.Marshal(message)
		if err != nil {
			log.Println(err)
			return
		}
		input.Data = bytes
	}

	inputByf, err := proto.Marshal(&input)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = c.Conn.Write(gim_util.Encode(inputByf))
	if err != nil {
		log.Println(err)
	}
}

func (c *TcpClient) Start() {
	connect, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		log.Println(err)
		return
	}

	c.codec = codecFactory.NewCodec(connect)
	c.Conn = connect

	c.SignIn()
	c.SyncTrigger()
	c.SubscribeRoom()
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
		bytes, err := c.codec.Read()
		if err != nil {
			log.Println(err)
			return
		}

		c.HandlePackage(bytes)
	}
}

func (c *TcpClient) SubscribeRoom() {
	c.Output(pb.PackageType_PT_SUBSCRIBE_ROOM, 0, &pb.SubscribeRoomInput{
		RoomId: 1,
		Seq:    0,
	})
}

func (c *TcpClient) HandlePackage(bytes []byte) {
	var output pb.Output
	err := proto.Unmarshal(bytes, &output)
	if err != nil {
		log.Println(err)
		return
	}

	switch output.Type {
	case pb.PackageType_PT_SIGN_IN:
		log.Println(Json(&output))
	case pb.PackageType_PT_HEARTBEAT:
		log.Println("心跳响应")
	case pb.PackageType_PT_SYNC:
		log.Println("离线消息同步开始------")
		syncResp := pb.SyncOutput{}
		err := proto.Unmarshal(output.Data, &syncResp)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("离线消息同步响应:code", output.Code, "message:", output.Message)
		for _, msg := range syncResp.Messages {
			log.Println(util.MessageToString(msg))
			c.Seq = msg.Seq
		}

		ack := pb.MessageACK{
			DeviceAck:   c.Seq,
			ReceiveTime: util.UnixMilliTime(time.Now()),
		}
		c.Output(pb.PackageType_PT_MESSAGE, output.RequestId, &ack)
		log.Println("离线消息同步结束------")
	case pb.PackageType_PT_MESSAGE:
		msg := pb.Message{}
		err := proto.Unmarshal(output.Data, &msg)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(util.MessageToString(&msg))
		c.Seq = msg.Seq
		ack := pb.MessageACK{
			DeviceAck:   msg.Seq,
			ReceiveTime: util.UnixMilliTime(time.Now()),
		}
		c.Output(pb.PackageType_PT_MESSAGE, output.RequestId, &ack)
	default:
		log.Println("switch other", output, len(bytes))
	}
}

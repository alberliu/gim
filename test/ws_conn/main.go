package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	"gim/pkg/protocol/pb"
	"gim/pkg/util"
)

func main() {
	client := WSClient{}
	fmt.Println("input UserId,DeviceId,SyncSequence")
	fmt.Scanf("%d %d %d", &client.UserId, &client.DeviceId, &client.Seq)
	client.Start()
	select {}
}

type WSClient struct {
	UserId   int64
	DeviceId int64
	Seq      int64
	Conn     *websocket.Conn
}

func (c *WSClient) Start() {
	conn, resp, err := websocket.DefaultDialer.Dial("ws://111.229.238.28:8002/ws", http.Header{})
	if err != nil {
		fmt.Println("dial error", err)
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bytes))
	c.Conn = conn

	c.SignIn()
	c.SyncTrigger()
	go c.Heartbeat()
	go c.Receive()
}

func (c *WSClient) Output(pt pb.PackageType, requestId int64, message proto.Message) {
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

	writeBytes, err := proto.Marshal(&input)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = c.Conn.WriteMessage(websocket.BinaryMessage, writeBytes)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *WSClient) SignIn() {
	signIn := pb.SignInInput{
		UserId:   c.UserId,
		DeviceId: c.DeviceId,
		Token:    "0",
	}
	c.Output(pb.PackageType_PT_SIGN_IN, time.Now().UnixNano(), &signIn)
}

func (c *WSClient) SyncTrigger() {
	c.Output(pb.PackageType_PT_SYNC, time.Now().UnixNano(), &pb.SyncInput{Seq: c.Seq})
}

func (c *WSClient) Heartbeat() {
	ticker := time.NewTicker(time.Minute * 4)
	for range ticker.C {
		c.Output(pb.PackageType_PT_HEARTBEAT, time.Now().UnixNano(), nil)
		fmt.Println("心跳发送")
	}
}

func (c *WSClient) Receive() {
	for {
		_, bytes, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		c.HandlePackage(bytes)
	}
}

func (c *WSClient) HandlePackage(bytes []byte) {
	var output pb.Output
	err := proto.Unmarshal(bytes, &output)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch output.Type {
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
		fmt.Printf("%+v \n", &output)
		for _, msg := range syncResp.Messages {
			log.Panicln(util.MessageToString(msg))
			c.Seq = msg.Seq
		}

		ack := pb.MessageACK{
			DeviceAck:   c.Seq,
			ReceiveTime: util.UnixMilliTime(time.Now()),
		}
		c.Output(pb.PackageType_PT_MESSAGE, output.RequestId, &ack)
		fmt.Println("离线消息同步结束------")
	case pb.PackageType_PT_MESSAGE:
		msg := pb.Message{}
		err := proto.Unmarshal(output.Data, &msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		log.Println("消息:", util.MessageToString(&msg))

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

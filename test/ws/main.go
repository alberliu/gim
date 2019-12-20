package main

import (
	"encoding/binary"
	"fmt"
	"gim/public/grpclib"
	"gim/public/pb"
	"gim/public/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/gorilla/websocket"

	"github.com/golang/protobuf/proto"
)

func main() {
	client := WSClient{}
	fmt.Println("input AppId,UserId,DeviceId,SyncSequence")
	fmt.Scanf("%d %d %d %d", &client.AppId, &client.UserId, &client.DeviceId, &client.Seq)
	client.Start()
	select {}
}

func Json(i interface{}) string {
	bytes, _ := jsoniter.Marshal(i)
	return string(bytes)
}

type WSClient struct {
	AppId    int64
	UserId   int64
	DeviceId int64
	Seq      int64
	Conn     *websocket.Conn
}

func (c *WSClient) Start() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}

	header := http.Header{}
	header.Set(grpclib.CtxAppId, strconv.FormatInt(c.AppId, 10))
	header.Set(grpclib.CtxUserId, strconv.FormatInt(c.UserId, 10))
	header.Set(grpclib.CtxDeviceId, strconv.FormatInt(c.DeviceId, 10))

	token, err := util.GetToken(c.AppId, c.UserId, c.DeviceId, time.Now().Add(24*30*time.Hour).Unix(), util.PublicKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	header.Set(grpclib.CtxToken, token)

	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		fmt.Println(err)
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bytes))
	c.Conn = conn

	c.SyncTrigger()
	go c.Heartbeat()
	go c.Receive()
}

func (c *WSClient) Output(pt int, message proto.Message) {
	var (
		bytes []byte
		err   error
	)
	if message != nil {
		bytes, err = proto.Marshal(message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	writeBytes := make([]byte, len(bytes)+2)
	binary.BigEndian.PutUint16(writeBytes[0:2], uint16(pt))
	copy(writeBytes[2:], bytes)
	err = c.Conn.WriteMessage(websocket.BinaryMessage, writeBytes)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *WSClient) SyncTrigger() {
	c.Output(int(pb.PackageType_PT_SYNC), &pb.SyncInput{Seq: c.Seq})
}

func (c *WSClient) Heartbeat() {
	ticker := time.NewTicker(time.Minute * 4)
	for range ticker.C {
		c.Output(int(pb.PackageType_PT_HEARTBEAT), nil)
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

		pt := int(binary.BigEndian.Uint16(bytes[0:2]))
		c.HandlePackage(pt, bytes[2:])
	}
}

func (c *WSClient) HandlePackage(pt int, bytes []byte) error {
	switch pb.PackageType(pt) {
	case pb.PackageType_PT_SIGN_IN:
		resp := pb.SignInOutput{}
		err := proto.Unmarshal(bytes, &resp)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(Json(resp))
	case pb.PackageType_PT_HEARTBEAT:
		fmt.Println("心跳响应")
	case pb.PackageType_PT_SYNC:
		fmt.Println("离线消息同步开始------")

		syncResp := pb.SyncOutput{}
		err := proto.Unmarshal(bytes, &syncResp)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("离线消息同步响应:code", syncResp.Code, "message:", syncResp.Message)
		for _, msg := range syncResp.Messages {
			if msg.ReceiverType == pb.ReceiverType_RT_USER {
				fmt.Printf("单聊消息：发送者类型：%d 发送者id：%d 接受者id：%d  消息内容：%+v seq：%d \n", msg.SenderType, msg.SenderId, msg.ReceiverId, msg.MessageBody.MessageContent, msg.Seq)
			}
			if msg.ReceiverType == pb.ReceiverType_RT_NORMAL_GROUP {
				fmt.Printf("群聊消息：发送者类型：%d 发送者id：%d 接受者id：%d  消息内容：%+v seq：%d \n", msg.SenderType, msg.SenderId, msg.ReceiverId, msg.MessageBody.MessageContent, msg.Seq)
			}
			if msg.ReceiverType == pb.ReceiverType_RT_LARGE_GROUP {
				fmt.Printf("大群消息：发送者类型：%d 发送者id：%d 接受者id：%d  消息内容：%+v seq：%d \n", msg.SenderType, msg.SenderId, msg.ReceiverId, msg.MessageBody.MessageContent, msg.Seq)
			}
		}
		fmt.Println("离线消息同步结束------")
	case pb.PackageType_PT_MESSAGE:
		message := pb.Message{}
		err := proto.Unmarshal(bytes, &message)
		if err != nil {
			fmt.Println(err)
			return err
		}

		msg := message.Message
		if msg.ReceiverType == pb.ReceiverType_RT_USER {
			fmt.Printf("单聊消息：发送者类型：%d 发送者id：%d 接受者id：%d  消息内容：%+v seq：%d \n", msg.SenderType, msg.SenderId, msg.ReceiverId, msg.MessageBody.MessageContent, msg.Seq)
		}
		if msg.ReceiverType == pb.ReceiverType_RT_NORMAL_GROUP {
			fmt.Printf("群聊消息：发送者类型：%d 发送者id：%d 接受者id：%d  消息内容：%+v seq：%d \n", msg.SenderType, msg.SenderId, msg.ReceiverId, msg.MessageBody.MessageContent, msg.Seq)
		}
		if msg.ReceiverType == pb.ReceiverType_RT_LARGE_GROUP {
			fmt.Printf("大群消息：发送者类型：%d 发送者id：%d 接受者id：%d  消息内容：%+v seq：%d \n", msg.SenderType, msg.SenderId, msg.ReceiverId, msg.MessageBody.MessageContent, msg.Seq)
		}

		c.Seq = msg.Seq
		c.Output(int(pb.PackageType_PT_MESSAGE), &pb.MessageACK{
			MessageId:   msg.MessageId,
			DeviceAck:   msg.Seq,
			ReceiveTime: util.UnixMilliTime(time.Now()),
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
	default:
		fmt.Println("switch other")
	}
	return nil
}

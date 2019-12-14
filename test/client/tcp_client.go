package client

import (
	"fmt"
	"gim/conn"
	"gim/public/logger"
	"gim/public/pb"
	"gim/public/util"
	"net"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/golang/protobuf/proto"
)

func Json(i interface{}) string {
	bytes, _ := jsoniter.Marshal(i)
	return string(bytes)
}

var codecFactory = conn.NewCodecFactory(2, 2, 65536, 1024)

type TcpClient struct {
	AppId    int64
	UserId   int64
	DeviceId int64
	Seq      int64
	codec    *conn.Codec
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
	go c.Receive()
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

	signInBytes, err := proto.Marshal(&signIn)
	if err != nil {
		fmt.Println(err)
		return
	}

	pack := conn.Package{Code: int(pb.PackageType_PT_SIGN_IN), Content: signInBytes}
	c.codec.Encode(pack, 10*time.Second)
}

func (c *TcpClient) SyncTrigger() {
	bytes, err := proto.Marshal(&pb.SyncInput{Seq: c.Seq})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.codec.Encode(conn.Package{Code: int(pb.PackageType_PT_SYNC), Content: bytes}, 10*time.Second)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *TcpClient) Heartbeat() {
	ticker := time.NewTicker(time.Minute * 4)
	for _ = range ticker.C {
		err := c.codec.Encode(conn.Package{Code: int(pb.PackageType_PT_HEARTBEAT), Content: []byte{}}, 10*time.Second)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("心跳发送")
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
			pack, ok, err := c.codec.Decode()
			if err != nil {
				fmt.Println(err)
				return
			}

			if ok {
				c.HandlePackage(*pack)
				continue
			}
			break
		}
	}
}

func (c *TcpClient) HandlePackage(pack conn.Package) error {
	switch pb.PackageType(pack.Code) {
	case pb.PackageType_PT_SIGN_IN:
		resp := pb.SignInOutput{}
		err := proto.Unmarshal(pack.Content, &resp)
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
		err := proto.Unmarshal(pack.Content, &syncResp)
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
		err := proto.Unmarshal(pack.Content, &message)
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

		ack := pb.MessageACK{
			MessageId:   msg.MessageId,
			DeviceAck:   msg.Seq,
			ReceiveTime: util.UnixMilliTime(time.Now()),
		}
		ackBytes, err := proto.Marshal(&ack)
		if err != nil {
			fmt.Println(err)
			return err
		}

		c.Seq = msg.Seq
		err = c.codec.Encode(conn.Package{Code: int(pb.PackageType_PT_MESSAGE), Content: ackBytes}, 10*time.Second)
		if err != nil {
			fmt.Println(err)
			return err
		}
	default:
		fmt.Println("switch other")
	}
	return nil
}

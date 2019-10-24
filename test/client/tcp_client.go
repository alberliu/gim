package client

import (
	"encoding/base64"
	"gim/connect"
	"gim/public/pb"
	"gim/public/util"
	"net"
	"strconv"
	"time"

	"github.com/json-iterator/go"

	"fmt"

	"github.com/golang/protobuf/proto"
)

func Json(i interface{}) string {
	bytes, _ := jsoniter.Marshal(i)
	return string(bytes)
}

type TcpClient struct {
	AppId    int64
	UserId   int64
	DeviceId int64
	Seq      int64
	codec    *connect.Codec
}

func (c *TcpClient) Start() {
	conn, err := net.Dial("tcp", "localhost:50000")
	if err != nil {
		fmt.Println(err)
		return
	}

	c.codec = connect.NewCodec(conn)

	c.SignIn()
	c.SyncTrigger()
	go c.HeadBeat()
	go c.Receive()
}

func (c *TcpClient) SignIn() {
	str := strconv.FormatInt(c.AppId, 10) + ":" + strconv.FormatInt(c.UserId, 10) + ":" +
		strconv.FormatInt(c.DeviceId, 10) + ":" + strconv.FormatInt(time.Now().Add(24*30*time.Hour).Unix(), 10)
	token, err := util.RsaEncrypt([]byte(str), []byte(util.PublicKey))
	if err != nil {
		fmt.Println(err)
		return
	}

	signIn := pb.SignInReq{
		AppId:    c.AppId,
		UserId:   c.UserId,
		DeviceId: c.DeviceId,
		Token:    base64.StdEncoding.EncodeToString(token),
	}

	signInBytes, err := proto.Marshal(&signIn)
	if err != nil {
		fmt.Println(err)
		return
	}

	pack := connect.Package{Code: int(pb.PackageType_PT_SIGN_IN_REQ), Content: signInBytes}
	c.codec.Encode(pack, 10*time.Second)
}

func (c *TcpClient) SyncTrigger() {
	bytes, err := proto.Marshal(&pb.SyncReq{Seq: c.Seq})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.codec.Encode(connect.Package{Code: int(pb.PackageType_PT_SYNC_REQ), Content: bytes}, 10*time.Second)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *TcpClient) HeadBeat() {
	ticker := time.NewTicker(time.Minute * 4)
	for _ = range ticker.C {
		err := c.codec.Encode(connect.Package{Code: int(pb.PackageType_PT_HEARTBEAT_REQ), Content: []byte{}}, 10*time.Second)
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

func (c *TcpClient) HandlePackage(pack connect.Package) error {
	switch pb.PackageType(pack.Code) {
	case pb.PackageType_PT_SIGN_IN_RESP:
		resp := pb.SignInResp{}
		err := proto.Unmarshal(pack.Content, &resp)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(Json(resp))
	case pb.PackageType_PT_HEARTBEAT_RESP:
		fmt.Println("心跳响应")
	case pb.PackageType_PT_SYNC_RESP:
		fmt.Println("离线消息同步开始------")

		syncResp := pb.SyncResp{}
		err := proto.Unmarshal(pack.Content, &syncResp)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("离线消息同步响应:code", syncResp.Code, "message:", syncResp.Code)
		for _, msg := range syncResp.Messages {
			if msg.ReceiverType == pb.ReceiverType_RT_USER {
				fmt.Println("单聊消息：发送者：", msg.SenderId, "接收者:", msg.SenderId, "内容:", msg.MessageBody.MessageContent.GetText())
			}
			if msg.ReceiverType == pb.ReceiverType_RT_NORMAL_GROUP {
				fmt.Println("小群消息：发送者：", msg.SenderId, "接收者:", msg.SenderId, "内容:", msg.MessageBody.MessageContent.GetText())
			}
			if msg.ReceiverType == pb.ReceiverType_RT_LARGE_GROUP {
				fmt.Println("大群消息：发送者：", msg.SenderId, "接收者:", msg.SenderId, "内容:", msg.MessageBody.MessageContent.GetText())
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
			fmt.Println("单聊消息：发送者：", msg.SenderId, "接收者:", msg.SenderId, "内容:", msg.MessageBody.MessageContent.GetText())
		}
		if msg.ReceiverType == pb.ReceiverType_RT_NORMAL_GROUP {
			fmt.Println("小群消息：发送者：", msg.SenderId, "接收者:", msg.SenderId, "内容:", msg.MessageBody.MessageContent.GetText())
		}
		if msg.ReceiverType == pb.ReceiverType_RT_LARGE_GROUP {
			fmt.Println("大群消息：发送者：", msg.SenderId, "接收者:", msg.SenderId, "内容:", msg.MessageBody.MessageContent.GetText())
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
		err = c.codec.Encode(connect.Package{Code: int(pb.PackageType_PT_MESSAGE_ACK), Content: ackBytes}, 10*time.Second)
		if err != nil {
			fmt.Println(err)
			return err
		}
	default:
		fmt.Println("switch other")
	}
	return nil
}

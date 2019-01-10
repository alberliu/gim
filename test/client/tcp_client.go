package client

import (
	"goim/connect"
	"goim/public/pb"
	"net"
	"time"

	"fmt"

	"goim/public/lib"

	"goim/public/transfer"

	"github.com/golang/protobuf/proto"
)

type TcpClient struct {
	DeviceId     int64
	UserId       int64
	Token        string
	SendSequence int64
	SyncSequence int64
	codec        *connect.Codec
}

func (c *TcpClient) Start() {
	conn, err := net.Dial("tcp", "localhost:50002")
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
	signIn := pb.SignIn{
		DeviceId: c.DeviceId,
		UserId:   c.UserId,
		Token:    c.Token,
	}

	signInBytes, err := proto.Marshal(&signIn)
	if err != nil {
		fmt.Println(err)
		return
	}

	pack := connect.Package{Code: connect.CodeSignIn, Content: signInBytes}
	c.codec.Eecode(pack, 10*time.Second)
}

func (c *TcpClient) SyncTrigger() {
	bytes, err := proto.Marshal(&pb.SyncTrigger{SyncSequence: c.SyncSequence})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.codec.Eecode(connect.Package{Code: connect.CodeSyncTrigger, Content: bytes}, 10*time.Second)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *TcpClient) HeadBeat() {
	ticker := time.NewTicker(time.Minute * 4)
	for _ = range ticker.C {
		err := c.codec.Eecode(connect.Package{Code: connect.CodeHeadbeat, Content: []byte{}}, 10*time.Second)
		if err != nil {
			fmt.Println(err)
		}
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
			pack, ok := c.codec.Decode()
			if ok {
				c.HandlePackage(*pack)
				continue
			}
			break
		}
	}
}

func (c *TcpClient) HandlePackage(pack connect.Package) error {
	fmt.Println("pack", pack.Code)
	switch pack.Code {
	case connect.CodeSignInACK:
		ack := pb.SignInACK{}
		err := proto.Unmarshal(pack.Content, &ack)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if ack.Code == 1 {
			fmt.Println("设备登录成功")
			return nil
		}
		fmt.Println("设备登录失败")

	case connect.CodeHeadbeatACK:
	case connect.CodeMessageSendACK:
		ack := pb.MessageSendACK{}
		err := proto.Unmarshal(pack.Content, &ack)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(ack.SendSequence, ack.Code)
	case connect.CodeMessage:
		message := pb.Message{}
		err := proto.Unmarshal(pack.Content, &message)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if message.Type == transfer.MessageTypeSync {
			fmt.Println("消息同步开始......")
		}

		for _, v := range message.Messages {
			if v.ReceiverType == 1 {
				if v.SenderDeviceId != c.DeviceId {
					fmt.Printf("单聊：来自用户：%d,消息内容：%s\n", v.SenderId, v.Content)
				}
			}
			if v.ReceiverType == 2 {
				if v.SenderDeviceId != c.DeviceId {
					fmt.Printf("群聊：来自用户：%d,群组：%d,消息内容：%s\n", v.SenderId, v.ReceiverId, v.Content)
				}
			}
			if c.SendSequence < v.SyncSequence {
				c.SendSequence = v.SyncSequence
			}
		}

		if message.Type == transfer.MessageTypeSync {
			fmt.Println("消息同步结束")
		}

		if len(message.Messages) == 0 {
			return nil
		}

		ack := pb.MessageACK{
			MessageId:    message.Messages[len(message.Messages)-1].MessageId,
			SyncSequence: message.Messages[len(message.Messages)-1].SyncSequence,
			ReceiveTime:  lib.UnixTime(time.Now()),
		}
		ackBytes, err := proto.Marshal(&ack)
		if err != nil {
			fmt.Println(err)
			return err
		}

		c.SyncSequence = ack.SyncSequence

		err = c.codec.Eecode(connect.Package{Code: connect.CodeMessageACK, Content: ackBytes}, 10*time.Second)
		if err != nil {
			fmt.Println(err)
			return err
		}
	default:
		fmt.Println("switch other")
	}
	return nil
}

func (c *TcpClient) SendMessage() {
	send := pb.MessageSend{}
	fmt.Scanf("%d %d %s", &send.ReceiverType, &send.ReceiverId, &send.Content)
	if send.Content == "" {
		return
	}
	send.Type = 1
	c.SendSequence++
	send.SendSequence = c.SendSequence
	send.SendTime = lib.UnixTime(time.Now())
	bytes, err := proto.Marshal(&send)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.codec.Eecode(connect.Package{Code: connect.CodeMessageSend, Content: bytes}, 10*time.Second)
	if err != nil {
		fmt.Println(err)
	}
}

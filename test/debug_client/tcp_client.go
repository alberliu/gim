package client

import (
	"goim/connect"
	"goim/public/pb"
	"net"
	"time"

	"goim/public/logger"

	"fmt"

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
		logger.Sugar.Error(err)
		return
	}

	c.codec = connect.NewCodec(conn)

	c.SignIn()
	c.SyncTrigger()
	//go c.HeadBeat()
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
	ticker := time.NewTicker(time.Second * 1)
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
			logger.Sugar.Error(err)
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
	switch pack.Code {
	case connect.CodeSignInACK:
		ack := pb.SignInACK{}
		err := proto.Unmarshal(pack.Content, &ack)
		if err != nil {
			logger.Sugar.Error(err)
			return err
		}
		logger.Sugar.Info("设备登录回执：%#v\n", ack)
	case connect.CodeHeadbeatACK:
		//log.Println("心跳回执")
	case connect.CodeMessageSendACK:
		ack := pb.MessageSendACK{}
		err := proto.Unmarshal(pack.Content, &ack)
		if err != nil {
			logger.Sugar.Error(err)
			return err
		}
		logger.Sugar.Info("消息发送回执：%#v\n", ack)
	case connect.CodeMessage:
		message := pb.Message{}
		err := proto.Unmarshal(pack.Content, &message)
		if err != nil {
			logger.Sugar.Error(err)
			return err
		}

		for _, v := range message.Messages {
			logger.Sugar.Info(message)
		}

		if len(message.Messages) == 0 {
			return nil
		}

		ack := pb.MessageACK{SyncSequence: message.Messages[len(message.Messages)-1].SyncSequence}
		ackBytes, err := proto.Marshal(&ack)
		if err != nil {
			logger.Sugar.Error(err)
			return err
		}

		c.SyncSequence = ack.SyncSequence

		err = c.codec.Eecode(connect.Package{Code: connect.CodeMessageACK, Content: ackBytes}, 10*time.Second)
		if err != nil {
			fmt.Println(err)
			return err
		}
	default:
		logger.Sugar.Info("switch other")
	}
	return nil
}

func (c *TcpClient) SendMessage() {
	send := pb.MessageSend{}
	fmt.Println("input ReceiverType,ReceiverId,Content")
	fmt.Scanf("%d %d %s", &send.ReceiverType, &send.ReceiverId, &send.Content)
	send.Type = 1
	c.SendSequence++
	send.SendSequence = c.SendSequence
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

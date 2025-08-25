package config

import (
	"log/slog"
	"os"

	"gim/pkg/protocol/pb/logicpb"
	"gim/pkg/protocol/pb/userpb"
)

const EnvLocal = "local"

var ENV = os.Getenv("ENV")

var builders = map[string]Builder{
	"local":   &localBuilder{},
	"compose": &composeBuilder{},
	"k8s":     &k8sBuilder{},
}

var Config Configuration

type Builder interface {
	Build() Configuration
}

type Configuration struct {
	LogLevel slog.Level
	LogFile  func(server string) string

	MySQL                string
	RedisHost            string
	RedisPassword        string
	PushRoomSubscribeNum int
	PushAllSubscribeNum  int

	ConnectLocalAddr     string
	ConnectRPCListenAddr string
	ConnectTCPListenAddr string
	ConnectWSListenAddr  string

	LogicRPCListenAddr string
	UserRPCListenAddr  string
	FileHTTPListenAddr string

	DeviceIntClientBuilder  func() logicpb.DeviceIntServiceClient
	MessageIntClientBuilder func() logicpb.MessageIntServiceClient
	RoomIntClientBuilder    func() logicpb.RoomIntServiceClient
	UserIntClientBuilder    func() userpb.UserIntServiceClient
}

func init() {
	builder, ok := builders[ENV]
	if !ok {
		builder = new(localBuilder)
	}
	Config = builder.Build()
}

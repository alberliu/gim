package config

import (
	"log/slog"
	"os"
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

	LogicRPCListenAddr    string
	BusinessRPCListenAddr string
	FileHTTPListenAddr    string

	LogicServerAddr    string
	BusinessServerAddr string
}

func init() {
	builder, ok := builders[ENV]
	if !ok {
		builder = new(localBuilder)
	}
	Config = builder.Build()
}

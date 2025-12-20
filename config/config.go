package config

import (
	"log/slog"
	"os"
)

const (
	EnvLocal   = "local"
	EnvCompose = "compose"
	EnvK8s     = "k8s"
)

var ENV = os.Getenv("ENV")

var builders = map[string]Builder{
	EnvLocal:   &localBuilder{},
	EnvCompose: &composeBuilder{},
	EnvK8s:     &k8sBuilder{},
}

const GrpcListenAddr = ":8000"

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
}

func init() {
	builder, ok := builders[ENV]
	if !ok {
		builder = new(localBuilder)
	}
	Config = builder.Build()
}

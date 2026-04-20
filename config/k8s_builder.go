package config

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/sercand/kuberesolver/v6"
	"google.golang.org/grpc/resolver"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gim/pkg/uk8s"
)

func init() {
	resolver.Register(kuberesolver.NewBuilder(nil, "k8s"))
}

type k8sBuilder struct{}

func (*k8sBuilder) Build() Configuration {
	const namespace = "default"

	k8sClient, err := uk8s.GetK8sClient()
	if err != nil {
		panic(err)
	}
	configmap, err := k8sClient.CoreV1().ConfigMaps(namespace).Get(context.TODO(), "config", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	return Configuration{
		LogLevel: slog.LevelDebug,
		LogFile: func(server string) string {
			return fmt.Sprintf("/data/log/%s/log.log", server)
		},

		MySQL:                configmap.Data["mysql"],
		RedisHost:            configmap.Data["redisIP"],
		RedisPassword:        configmap.Data["redisPassword"],
		PushRoomSubscribeNum: getInt(configmap.Data, "pushRoomSubscribeNum"),
		PushAllSubscribeNum:  getInt(configmap.Data, "pushAllSubscribeNum"),
	}
}

func getInt(m map[string]string, key string) int {
	value, _ := strconv.Atoi(m[key])
	return value
}

package config

import (
	"context"
	"gim/pkg/k8sutil"
	"gim/pkg/logger"
	"os"
	"strconv"

	"go.uber.org/zap"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	RPCListenAddr = ":8000"
	TCPListenAddr = ":8080"
	WSListenAddr  = ":8001"
)

var (
	Namespace     = "gimns"
	MySQL         string
	RedisIP       string
	RedisPassword string

	LocalAddr            string
	PushRoomSubscribeNum int
	PushAllSubscribeNum  int
)

func Init() {
	k8sClient, err := k8sutil.GetK8sClient()
	if err != nil {
		panic(err)
	}
	configmap, err := k8sClient.CoreV1().ConfigMaps(Namespace).Get(context.TODO(), "config", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	MySQL = configmap.Data["mysql"]
	RedisIP = configmap.Data["redisIP"]
	RedisPassword = configmap.Data["redisPassword"]
	PushRoomSubscribeNum, _ = strconv.Atoi(configmap.Data["pushRoomSubscribeNum"])
	if PushRoomSubscribeNum == 0 {
		panic("PushRoomSubscribeNum == 0")
	}
	PushAllSubscribeNum, _ = strconv.Atoi(configmap.Data["pushAllSubscribeNum"])
	if PushRoomSubscribeNum == 0 {
		panic("PushAllSubscribeNum == 0")
	}

	LocalAddr = os.Getenv("POD_IP") + RPCListenAddr

	logger.Level = zap.DebugLevel
	logger.Target = logger.Console
}

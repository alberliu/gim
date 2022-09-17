package config

import (
	"context"
	"gim/pkg/k8sutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	RPCListenAddr = ":8000"
	TCPListenAddr = ":8080"
	WSListenAddr  = ":8001"
)

var (
	NameSpace     string
	MySQL         string
	RedisIP       string
	RedisPassword string

	LocalAddr            string
	PushRoomSubscribeNum int
	PushAllSubscribeNum  int
)

func init() {
	k8sClient, err := k8sutil.GetK8sClient()
	if err != nil {
		panic(err)
	}

	configmap, err := k8sClient.CoreV1().ConfigMaps(NameSpace).Get(context.TODO(), "config", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	MySQL = configmap.Data["mysql"]
	RedisIP = configmap.Data["redisIP"]
	RedisPassword = configmap.Data["redisPassword"]
	PushRoomSubscribeNum = configmap.Data["pushRoomSubscribeNum"]
	PushAllSubscribeNum = configmap.Data["pushAllSubscribeNum"]
}

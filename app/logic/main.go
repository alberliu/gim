package main

import (
	"fmt"
	"goim/conf"
	"goim/logic/controller"
	"goim/logic/mq/consume"
	"goim/logic/rpc/server"
	"goim/public/logger"
	"runtime"
)

func main() {
	// 启动rpc服务
	go func() {
		defer RecoverPanic()
		server.StartRPCServer()
	}()

	// 启动nsq消费服务
	go func() {
		defer RecoverPanic()
		consume.StartNsqConsumer()
	}()

	// 启动web容器
	controller.Engine.Run(conf.LogicHTTPListenIP)
}

// RecoverPanic 恢复panic
func RecoverPanic() {
	err := recover()
	if err != nil {
		logger.Sugar.Error(err)
		logger.Sugar.Error(GetPanicInfo())
	}
}

// PrintStaStack 打印Panic堆栈信息
func GetPanicInfo() string {
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	return fmt.Sprintf("%s", buf[:n])
}

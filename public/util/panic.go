package util

import (
	"fmt"
	"goim/public/logger"
	"runtime"
)

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

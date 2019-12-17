package imerror

import (
	"fmt"
	"runtime"
	"strings"
)

const name = "gim"

// UnknownError 位置错误
type UnknownError struct {
	Err   error
	Stack []string
}

func (e *UnknownError) Error() string {
	return e.Err.Error()
}

func WrapError(err error) error {
	e := &UnknownError{
		Err:   err,
		Stack: stack(),
	}
	return e
}

// Stack 获取堆栈信息
func stack() []string {
	var pc = make([]uintptr, 20)
	n := runtime.Callers(3, pc)

	stack := make([]string, 0, n)
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i] - 1)
		file, line := f.FileLine(pc[i] - 1)
		n := strings.Index(file, name)
		if n != -1 {
			s := fmt.Sprintf("%s:%d %s", file, line, f.Name())
			stack = append(stack, s)
		}
	}
	return stack
}

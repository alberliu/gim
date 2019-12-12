package imerror

import (
	"fmt"
	"runtime"
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
		Stack: Stack(),
	}
	return e
}

// Stack 获取堆栈信息
func Stack() []string {
	var pc = make([]uintptr, 20)
	n := runtime.Callers(3, pc)

	stack := make([]string, 0, n)
	/*for i := 0; i < n; i++ {
		//f := runtime.FuncForPC(pc[i])
		//file, line := f.FileLine(pc[i])
		//n := strings.Index(file, name)
		//if n != -1 {
		//s := fmt.Sprintf("%s:%d %s", file, line, f.Name())
		//stack = append(stack, s)
		//}

		stack = append(stack, fmt.Sprintf("%+v", pc[i]))
	}*/

	fmt.Printf("%+v", pc[:n])
	return stack
}

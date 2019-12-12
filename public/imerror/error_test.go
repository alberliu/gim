package imerror

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	err := a()
	e := err.(*UnknownError)
	fmt.Println(e.Err)
	for _, v := range e.Stack {
		fmt.Println(v)
	}
}

func a() error {
	fmt.Println("11111111111111111")
	return b()
}

func b() error {
	fmt.Println("11111111111111111")
	return c()
}
func c() error {
	fmt.Println("11111111111111111")
	fmt.Println("11111111111111111")
	fmt.Println("11111111111111111")
	fmt.Println("11111111111111111")
	fmt.Println("11111111111111111")
	fmt.Println("11111111111111111")
	fmt.Println("11111111111111111")
	return WrapError(errors.New("kkk"))
}

package imerror

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	err := a()
	e := err.(*UnknownError)
	for _, v := range e.Stack {
		fmt.Println(v)
	}
}

func a() error {
	return b()
}

func b() error {
	return c()
}
func c() error {
	return WrapError(errors.New("kkk"))
}

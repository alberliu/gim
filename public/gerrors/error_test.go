package gerrors

import (
	"errors"
	"fmt"
	"testing"

	"google.golang.org/grpc/status"
)

func TestError(t *testing.T) {
	s, ok := status.FromError(errors.New("err"))
	fmt.Println(ok)
	fmt.Printf("%+v", *s)
	fmt.Println(s.Code())
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

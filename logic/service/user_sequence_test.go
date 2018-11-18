package service

import (
	"fmt"
	"testing"
)

func TestUserRequenceService_GetNext(t *testing.T) {
	fmt.Println(UserRequenceService.GetNext(ctx, 1))
}

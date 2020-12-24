package cache

import (
	"fmt"
	"gim/internal/business/model"
	"testing"
)

func TestAuthCache_Get(t *testing.T) {
	fmt.Println(AuthCache.Get(1, 1))
}

func TestAuthCache_Set(t *testing.T) {
	fmt.Println(AuthCache.Set(1, 1, model.Device{
		Type:   1,
		Token:  "111",
		Expire: 111,
	}))
}

func TestAuthCache_GetAll(t *testing.T) {
	fmt.Println(AuthCache.GetAll(1))
}

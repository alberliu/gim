package util

import (
	"fmt"
	"testing"
)

// appId:user_id:device_id:expire token格式
func TestRsaEncrypt(t *testing.T) {
	token, err := GetToken(1000000000000000, 100000000000000000, 100000000000000000, 1000000000000000, PublicKey)
	fmt.Println(err)
	fmt.Println(token)
	fmt.Println(len(token))

	fmt.Println(DecryptToken(token, PrivateKey))
}

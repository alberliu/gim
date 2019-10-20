package util

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"testing"
	"time"
)

// appId:user_id:device_id:expire token格式
func TestRsaEncrypt(t *testing.T) {
	str := "1:1:1:" + strconv.FormatInt(time.Now().Add(24*30*time.Hour).Unix(), 10)
	token, err := RsaEncrypt([]byte(str), PublicKey)
	fmt.Println(err)
	fmt.Println(base64.StdEncoding.EncodeToString(token))
	fmt.Println(string(token))
}

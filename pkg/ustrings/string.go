package ustrings

import (
	"crypto/rand"
	"math/big"
)

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// RandString 生成随机字符串
func RandString(n int) string {
	bytes := make([]byte, n)
	for i := range bytes {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		bytes[i] = letters[idx.Int64()]
	}
	return string(bytes)
}

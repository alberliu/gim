package ustrings

import (
	"crypto/rand"
	"math/big"
)

// RandString 生成随机字符串
func RandString(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	bytes := make([]byte, n)
	for i := range bytes {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		bytes[i] = letters[idx.Int64()]
	}
	return string(bytes)
}

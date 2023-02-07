package util

import (
	"fmt"
)

func FormatMessage(code int32, content []byte) string {
	return fmt.Sprintf("code:%d,content:%s", code, string(content))
}

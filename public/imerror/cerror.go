package imerror

var (
	CCodeSuccess = 0 // 成功发送
)

// CError 接入层调用错误
type CError struct {
	Code    int
	Message string
}

func (e *CError) Error() string {
	return e.Message
}

func NewCError(code int, message string) *CError {
	return &CError{
		Code:    code,
		Message: message,
	}
}

var (
	CErrUnkonw      = NewCError(1, "unknow error")  // 未知错误
	CErrNotIsFriend = NewCError(2, "not is friend") // 非好友关系
	CErrNotInGroup  = NewCError(3, "not in group")  // 没有在群组内
)

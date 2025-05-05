package gerrors

var (
	ErrBadCode      = newError(10300, "验证码错误")
	ErrUserNotFound = newError(10301, "用户找不到")
)

package gerrors

var (
	ErrConnNotFound         = newError(10101, "连接找不到")
	ErrConnDeviceIdNotEqual = newError(10102, "连接设备ID不相等")
)

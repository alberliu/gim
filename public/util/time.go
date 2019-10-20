package util

import "time"

// FormatTime 格式化时间
func FormatTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

// ParseTime 将时间字符串转为Time
func ParseTime(str string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", str)
}

// UnixMilliTime 将时间转化为毫秒数
func UnixMilliTime(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

// UnunixMilliTime 将毫秒数转为为时间
func UnunixMilliTime(unix int64) time.Time {
	return time.Unix(0, unix*1000000)
}

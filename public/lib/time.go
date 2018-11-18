package lib

import "time"

// FormatTime 格式化时间
func FormatTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

// UnFormatTime 将时间字符串转为Time
func UnFormatTime(str string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", str)
}

// UnixTime 将时间转化为毫秒数
func UnixTime(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

// UnunixTime 将毫秒数转为为时间
func UnunixTime(unix int64) time.Time {
	return time.Unix(0, unix*1000000)
}

package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"

	"gim/config"
)

func Init(directory string) {
	var writer io.Writer

	logFile := config.Config.LogFile(directory)
	if logFile == "" {
		writer = os.Stdout
	} else {
		writer = &lumberjack.Logger{
			Filename:   fmt.Sprintf("/data/log/%s/log.log", directory),
			MaxSize:    100, // 单个文件大小megabytes
			MaxBackups: 30,  // 最大备份数量
			MaxAge:     30,  // 保存天数
			LocalTime:  true,
		}
	}

	options := &slog.HandlerOptions{
		AddSource:   true,
		Level:       config.Config.LogLevel,
		ReplaceAttr: replaceAttr,
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(writer, options)))
	slog.Info("slog init")
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case "time":
		a.Key = "ts"
		a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05.000"))
	case "level":
		a.Value = slog.StringValue(strings.ToLower(a.Value.String()))
	case "source":
		source := a.Value.Any().(*slog.Source)
		a.Value = slog.StringValue(getShortSource(source))
	}
	return a
}

func getShortSource(source *slog.Source) string {
	index := strings.LastIndex(source.File, "/")
	if index != -1 {
		index = strings.LastIndex(source.File[0:index], "/")
	}
	return strings.ReplaceAll(source.File[index+1:], ".go", "") + ":" + strconv.Itoa(source.Line)
}

func Error(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

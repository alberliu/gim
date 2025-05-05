package gerrors

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"gim/pkg/util"
)

const name = "im"

const TypeUrlStack = "type_url_stack"

func GetErrorStack(s *status.Status) string {
	pbs := s.Proto()
	for i := range pbs.Details {
		if pbs.Details[i].TypeUrl == TypeUrlStack {
			return util.Bytes2str(pbs.Details[i].Value)
		}
	}
	return ""
}

// Stack 获取堆栈信息
func stack() string {
	var pc = make([]uintptr, 20)
	n := runtime.Callers(3, pc)

	var build strings.Builder
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i] - 1)
		file, line := f.FileLine(pc[i] - 1)
		n := strings.Index(file, name)
		if n != -1 {
			s := fmt.Sprintf(" %s:%d \n", file[n:], line)
			build.WriteString(s)
		}
	}
	return build.String()
}

func LogPanic(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, err *error) {
	p := recover()
	if p != nil {
		slog.Error("panic", "info", info, "ctx", ctx, "req", req, "panic", p,
			"stack", util.GetStackInfo())
		*err = ErrUnknown
	}
}

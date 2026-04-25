package gerrors

import (
	"context"
	"log/slog"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const TypeUrlStack = "type_url_stack"

func GetErrorStack(s *status.Status) string {
	pbs := s.Proto()
	for i := range pbs.Details {
		if pbs.Details[i].TypeUrl == TypeUrlStack {
			return string(pbs.Details[i].Value)
		}
	}
	return ""
}

func LogPanic(ctx context.Context, req any, info *grpc.UnaryServerInfo, err *error) {
	p := recover()
	if p != nil {
		slog.Error("panic recovered", "info", info, "ctx", ctx, "req", req, "panic", p,
			"stack", string(debug.Stack()))
		*err = ErrUnknown
	}
}

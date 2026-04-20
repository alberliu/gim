package gerrors

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"gim/pkg/util"
)

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

func LogPanic(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, err *error) {
	p := recover()
	if p != nil {
		slog.Error("panic", "info", info, "ctx", ctx, "req", req, "panic", p,
			"stack", util.GetStackInfo())
		*err = ErrUnknown
	}
}

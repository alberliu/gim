package logger

import (
	"context"
	"log/slog"

	"gim/pkg/md"
)

type ctxHandler struct {
	inner slog.Handler
}

func NewHandler(inner slog.Handler) slog.Handler {
	return &ctxHandler{
		inner: inner,
	}
}

func (h *ctxHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *ctxHandler) Handle(ctx context.Context, r slog.Record) error {
	requestID := md.GetRequestID(ctx)
	if requestID != "" {
		r.AddAttrs(slog.String("request_id", requestID))
	}
	return h.inner.Handle(ctx, r)
}

func (h *ctxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ctxHandler{inner: h.inner.WithAttrs(attrs)}
}

func (h *ctxHandler) WithGroup(name string) slog.Handler {
	return &ctxHandler{inner: h.inner.WithGroup(name)}
}
